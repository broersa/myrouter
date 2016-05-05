package dalpsql

import (
	"database/sql"

	"github.com/broersa/mylog"
	"github.com/broersa/myrouter/dal"
)

type (
	dalpsql struct {
		tx *sql.Tx
	}
)

// NewDal ...
func NewDal(tx *sql.Tx) dal.Dal {
	mylog.Trace.Println("DalPsql created.")
	return &dalpsql{tx}
}

// Commit ...
func (d *dalpsql) Commit() error {
	err := d.tx.Commit()
	return err
}

// Rollback ...
func (d *dalpsql) Rollback() error {
	err := d.tx.Rollback()
	return err
}

// AddBroker ...
func (d *dalpsql) AddBroker(broker *dal.Broker) (int64, error) {
	q, err := d.tx.Prepare("insert into brokers (broname, broendpoint) values ($1, $2) returning brokey")
	if err != nil {
		return 0, err
	}
	var r int64
	err = d.tx.Stmt(q).QueryRow(broker.Name, broker.Endpoint).Scan(&r)
	if err != nil {
		return 0, err
	}
	return r, nil
}

// GetBrokers ...
func (d *dalpsql) GetBrokers() ([]dal.Broker, error) {
	var returnvalue []dal.Broker
	rows, err := d.tx.Query("SELECT brokey, broname, broendpoint FROM brokers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		br := dal.Broker{}
		if err := rows.Scan(&br.ID, &br.Name, &br.Endpoint); err != nil {
			return nil, err
		}
		returnvalue = append(returnvalue, br)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return returnvalue, nil
}
