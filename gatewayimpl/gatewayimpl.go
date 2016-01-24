package gatewayimpl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/broersa/semtech"
	"github.com/broersa/ttnrouter/gateway"
)

type (
	GatewayImpl struct {
	}
)

// New Implemented Factory
func New() gateway.Gateway {
	return &GatewayImpl{}
}

func (gatewayimpl *GatewayImpl) SendPushACK(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16) error {
	out := make([]byte, 4)
	out[0] = protocolversion
	binary.LittleEndian.PutUint16(out[1:3], randomtoken)
	out[3] = byte(semtech.PushACK)
	_, err := conn.WriteToUDP(out, addr)
	if err != nil {
		return err
	}
	return nil
}

func (gatewayimpl *GatewayImpl) SendPullACK(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16) error {
	out := make([]byte, 4)
	out[0] = protocolversion
	binary.LittleEndian.PutUint16(out[1:3], randomtoken)
	out[3] = byte(semtech.PullACK)
	_, err := conn.WriteToUDP(out, addr)
	if err != nil {
		return err
	}
	return nil
}

func (gatewayimpl *GatewayImpl) SendPullResp(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16, txpk string) error {
	json := "{\"txpk\":" + txpk + "}"

	b0 := new(bytes.Buffer)
	b0.WriteByte(protocolversion)
	binary.Write(b0, binary.LittleEndian, randomtoken)
	b0.WriteByte(byte(semtech.PullResp))
	b0.Write([]byte(json))
	_, err := conn.WriteToUDP(b0.Bytes(), addr)
	if err != nil {
		return err
	}
	fmt.Println(string(b0.Bytes()))
	return nil
}
