package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/broersa/lora"
	"github.com/broersa/myrouter/bllimpl"
	"github.com/broersa/myrouter/broker"
	"github.com/broersa/myrouter/brokerimpl"
	"github.com/broersa/myrouter/dalpsql"
	"github.com/broersa/myrouter/gatewayimpl"
	"github.com/broersa/semtech"

	// Database Driver

	_ "github.com/lib/pq"
)

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Print("MYRouter is ALIVE")
	c := os.Getenv("MYROUTER_DB")
	//s, err := sql.Open("postgres", "postgres://user:password@server/db?sslmode=require")
	s, err := sql.Open("postgres", c)
	checkerror(err)
	d := dalpsql.New(s)
	bll := bllimpl.New(&d)
	gw := gatewayimpl.New()
	bro := brokerimpl.New()
	brokers, err := bll.GetBrokers()
	checkerror(err)
	brokerlist := make([]string, 0)
	for _, value := range brokers {
		brokerlist = append(brokerlist, value.Endpoint)
	}

	ServerAddr, err := net.ResolveUDPAddr("udp", ":1700")
	checkerror(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	checkerror(err)

	defer ServerConn.Close()

	buf := make([]byte, 2048)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		checkerror(err)
		packettype, err := semtech.GetPacketType(buf)
		checkerror(err)
		switch packettype {
		case semtech.PushData:
			pd := new(semtech.PushDataPacket)
			err := pd.UnmarshalBinary(buf[0:n])
			checkerror(err)
			err = gw.SendPushACK(ServerConn, addr, pd.ProtocolVersion, pd.RandomToken)
			checkerror(err)
			fmt.Println("<-PushACK")
			for _, value := range pd.Payload.RXPK {
				data, err := base64.StdEncoding.DecodeString(value.Data)
				checkerror(err)
				mhdr, err := lora.NewMHDRFromByte(data[0])
				if err != nil {
					if _, ok := err.(*lora.ErrorMTypeValidationFailed); ok {
						log.Print(err)
					} else {
						if _, ok := err.(*lora.ErrorMajorValidationFailed); ok {
							log.Print(err)
						} else {
							checkerror(err)
						}
					}
				} else {
					if mhdr.IsJoinRequest() {
						fmt.Println("->JoinRequest")
						joinrequest, err := lora.NewJoinRequest(data)
						checkerror(err)
						fmt.Println(hex.EncodeToString(joinrequest.GetAppEUI()))
						fmt.Println(hex.EncodeToString(joinrequest.GetDevEUI()))
						message := &broker.Message{addr.Network(), addr.String(), ServerAddr.Network(), ServerAddr.String(), value}
						endpoint, err := bro.FindBrokerOnAppEUI(joinrequest.GetAppEUI(), brokerlist)
						if err != nil {
							log.Print(err)
						} else {
							responsemessage, err := bro.ForwardMessage(endpoint, message)
							checkerror(err)
							txpk, err := json.Marshal(responsemessage.Package)
							gateway, err := bll.FindGateway(pd.GatewayMAC[:])
							checkerror(err)
							err = gw.SendPullResp(ServerConn, gateway, pd.ProtocolVersion, pd.RandomToken, string(txpk))
							checkerror(err)
							fmt.Println("<-JoinAccept")
						}
					}
				}
			}
		case semtech.PullData:
			pd := new(semtech.PullDataPacket)
			err := pd.UnmarshalBinary(buf[0:n])
			checkerror(err)
			fmt.Println("->PullData")
			bll.RefreshGateway(pd.GatewayMAC[:], addr)
			err = gw.SendPullACK(ServerConn, addr, pd.ProtocolVersion, pd.RandomToken)
			checkerror(err)
			fmt.Println("<-PullACK")
		}
	}
}
