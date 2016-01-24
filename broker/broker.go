package broker

type (
	// Broker Interface for Business Layer
	Broker interface {
		FindBrokerOnAppEUI(appeui []byte, brokers []string) (string, error)
		FindBrokerOnDevAddr(devaddr []byte, brokers []string) (string, error)
		ForwardMessage(endpoint string, message *Message) (*ResponseMessage, error)
	}
)
