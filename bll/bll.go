package bll

import "net"

type (
	// Bll Interface for Business Layer
	Bll interface {
		GetBrokers() ([]Broker, error)
		RefreshGateway(gatewaymac []byte, addr *net.UDPAddr) error
		FindGateway(gatewaymac []byte) (*net.UDPAddr, error)
	}
)
