package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/broersa/lora"
	"github.com/broersa/semtech"
	"github.com/broersa/ttnrouter/bllimpl"
	"github.com/broersa/ttnrouter/broker"
	"github.com/broersa/ttnrouter/brokerimpl"
	"github.com/broersa/ttnrouter/dalpsql"
	"github.com/broersa/ttnrouter/gatewayimpl"

	// Database Driver

	_ "github.com/lib/pq"
)

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Print("TTNRouter is ALIVE")
	c := os.Getenv("TTNROUTER_DB")
	//s, err := sql.Open("postgres", "postgres://user:password@server/ttn?sslmode=require")
	s, err := sql.Open("postgres", c)
	checkerror(err)

	d := dalpsql.New(s)
	//err = d.BeginTransaction()
	//checkerror(err)
	//b, err := d.GetBrokers() //(&dal.Broker{Name: "tester", Endpoint: "http://127.0.0.1:4333"})
	//checkerror(err)
	//err = d.CommitTransaction()
	//checkerror(err)

	//for _, value := range b {
	//	fmt.Println(value.Name)
	//}
	bll := bllimpl.New(&d)
	gw := gatewayimpl.New()
	bro := brokerimpl.New()
	l, err := bll.GetBrokers()
	checkerror(err)
	ll := make([]string, 0)
	for _, value := range l {
		ll = append(ll, value.Endpoint)
	}

	appx := make([]byte, 0)
	e, err := bro.FindBrokerOnAppEUI(appx, ll)
	checkerror(err)
	log.Fatal(e)
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
						message := &broker.Message{addr.Network(), addr.String(), ServerAddr.Network(), ServerAddr.String(), value}
						appeui, err := bro.FindBrokerOnAppEUI(joinrequest.GetAppEUI(), make([]string, 0))
						checkerror(err)
						responsemessage, err := bro.ForwardMessage(appeui, message)
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