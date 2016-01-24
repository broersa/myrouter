package dalpsql

import (
	"database/sql"

	"github.com/broersa/myrouter/dal"
)

type (
	dalPsql struct {
		db *sql.DB
		tx *sql.Tx
	}
)

// New Implemented Factory
func New(db *sql.DB) dal.Dal {
	return &dalPsql{db, nil}
}

func (dalpsql *dalPsql) BeginTransaction() error {
	tx, err := dalpsql.db.Begin()
	if err != nil {
		return err
	}
	dalpsql.tx = tx
	return nil
}

func (dalpsql *dalPsql) CommitTransaction() error {
	err := dalpsql.tx.Commit()
	return err
}

func (dalpsql *dalPsql) RollbackTransaction() error {
	err := dalpsql.tx.Rollback()
	return err
}

func (dalpsql *dalPsql) AddBroker(broker *dal.Broker) (int64, error) {
	q, err := dalpsql.db.Prepare("insert into brokers (broname, broendpoint) values ($1, $2) returning brokey")
	if err != nil {
		return 0, err
	}
	var r int64
	err = dalpsql.tx.Stmt(q).QueryRow(broker.Name, broker.Endpoint).Scan(&r)
	if err != nil {
		return 0, err
	}
	return r, nil
}

func (dalpsql *dalPsql) GetBrokers() ([]dal.Broker, error) {
	var returnvalue []dal.Broker
	rows, err := dalpsql.db.Query("SELECT brokey, broname, broendpoint FROM brokers")
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
