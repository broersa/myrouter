package dal

type (
	// Dal Interface for Business Layer
	Dal interface {
		Commit() error
		Rollback() error
		AddBroker(broker *Broker) (int64, error)
		GetBrokers() ([]Broker, error)
		//		GetSessionOnID(id int64) (*Session, error)
		//		GetSessionOnDevNonce(devnonce string) (*Session, error)
	}
)
