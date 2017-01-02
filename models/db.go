package models

import (
	"io/ioutil"
	"math/big"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Datastore interface {
	InsertUser(string, []byte, string) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserById(string) (*User, error)

	GetLogsByUser(string) (*[]Logs, error)
	InsertLog(string) (*Logs, error)
	UpdateLog(string) (*Logs, error)
	//DeleteLog(string) error

	GetEntriesByUser(string) (*[]Entry, error)
	GetEntriedByLog(string) (*[]Entry, error)
	InsertEntry(string, *big.Int, float64, float64, string) (*Entry, error)
	//DeleteEntry(string) error
}

type Data struct {
	*sqlx.DB
}

func CreateDB(database string, connectionOptions string) (*Data, error) {
	// Import database file and convert to string
	schema, err := ioutil.ReadFile(database)
	if err != nil {
		return nil, err
	}
	sql := string(schema[:])

	db, err := sqlx.Connect("postgres", connectionOptions)
	if err != nil {
		return nil, err
	}

	db.MustExec(sql)

	return &Data{db}, nil
}
