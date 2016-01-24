package bllimpl

import (
	"net"

	"github.com/broersa/ttnrouter/bll"
	"github.com/broersa/ttnrouter/dal"
)

type (
	BllImpl struct {
		dal      *dal.Dal
		gateways map[string]*net.UDPAddr
	}
)

// New Implemented Factory
func New(dal *dal.Dal) bll.Bll {
	return &BllImpl{dal: dal, gateways: make(map[string]*net.UDPAddr)}
}

func (bllimpl *BllImpl) GetBrokers() ([]bll.Broker, error) {
	returnvalue := make([]bll.Broker, 0)
	b, err := (*bllimpl.dal).GetBrokers()
	if err != nil {
		return nil, err
	}
	for _, value := range b {
		returnvalue = append(returnvalue, bll.Broker{ID: value.ID, Name: value.Name, Endpoint: value.Endpoint})
	}
	return returnvalue, nil
}

func (bllimpl *BllImpl) RefreshGateway(gatewaymac []byte, from *net.UDPAddr) error {
	bllimpl.gateways[string(gatewaymac)] = from
	return nil
}

func (bllimpl *BllImpl) FindGateway(gatewaymac []byte) (*net.UDPAddr, error) {
	return bllimpl.gateways[string(gatewaymac)], nil
}
