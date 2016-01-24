package broker

import (
	"github.com/broersa/semtech"
)

// Message Data Entity
type Message struct {
	OriginUDPAddrNetwork string       `json:"originudpaddrnetwork"`
	OriginUDPAddrString  string       `json:"originudpaddrstring"`
	ReturnUDPAddrNetwork string       `json:"returnudpaddrnetwork"`
	ReturnUDPAddrString  string       `json:"returnudpaddrstring"`
	Package              semtech.RXPK `json:"package"`
}
