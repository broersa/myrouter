package dal

type (
	// Dal Interface for Business Layer
	Dal interface {
		BeginTransaction() error
		CommitTransaction() error
		RollbackTransaction() error
		AddBroker(broker *Broker) (int64, error)
		GetBrokers() ([]Broker, error)
		//		GetSessionOnID(id int64) (*Session, error)
		//		GetSessionOnDevNonce(devnonce string) (*Session, error)
	}
)
