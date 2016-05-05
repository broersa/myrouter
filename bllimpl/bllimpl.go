package bllimpl

import (
	"net"

	"github.com/broersa/myrouter/bll"
	"github.com/broersa/myrouter/dal"
)

type (
	bllimpl struct {
		dal      *dal.Factory
		gateways map[string]*net.UDPAddr
	}
)

// NewBll Implemented Factory
func NewBll(dal *dal.Factory) bll.Bll {
	return &bllimpl{dal: dal, gateways: make(map[string]*net.UDPAddr)}
}

// GetBrokers ...
func (b *bllimpl) GetBrokers() ([]bll.Broker, error) {
	var returnvalue []bll.Broker
	tx, err := (*b.dal).GetInstance()
	if err != nil {
		return nil, err
	}
	br, err := tx.GetBrokers()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, value := range br {
		returnvalue = append(returnvalue, bll.Broker{ID: value.ID, Name: value.Name, Endpoint: value.Endpoint})
	}
	tx.Commit()
	return returnvalue, nil
}

// RefreshGateway ...
func (b *bllimpl) RefreshGateway(gatewaymac []byte, from *net.UDPAddr) error {
	b.gateways[string(gatewaymac)] = from
	return nil
}

// FindGateway ...
func (b *bllimpl) FindGateway(gatewaymac []byte) (*net.UDPAddr, error) {
	return b.gateways[string(gatewaymac)], nil
}
