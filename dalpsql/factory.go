package dalpsql

import (
	"database/sql"

	"github.com/broersa/mylog"
	"github.com/broersa/myrouter/dal"
)

type (
	// Factory ...
	factory struct {
		db *sql.DB
	}
)

// NewFactory ...
func NewFactory(db *sql.DB) dal.Factory {
	mylog.Trace.Println("Factory created.")
	return &factory{db}
}

// NewFactoryInstance ...
func (f *factory) GetInstance() (dal.Dal, error) {
	tx, err := f.db.Begin()
	if err != nil {
		return nil, err
	}
	return NewDal(tx), nil
}
