package broker

import (
	"github.com/broersa/semtech"
)

// Message Data Entity
type ResponseMessage struct {
	OriginUDPAddrNetwork string       `json:"originudpaddrnetwork"`
	OriginUDPAddrString  string       `json:"originudpaddrstring"`
	Package              semtech.TXPK `json:"package"`
}
