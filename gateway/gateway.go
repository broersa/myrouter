package gateway

import "net"

type (
	// Gateway Interface for Business Layer
	Gateway interface {
		SendPushACK(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16) error
		SendPullACK(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16) error
		SendPullResp(conn *net.UDPConn, addr *net.UDPAddr, protocolversion byte, randomtoken uint16, txpk string) error
	}
)
