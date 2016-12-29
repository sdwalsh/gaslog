package models

import (
	"math/big"
	"time"
)

type Entry struct {
	Id              string    `db:"id"`
	LogBook         string    `db:"log_book_id"`
	Created_at      time.Time `db:"created_at"`
	Miles           int64  		`db:"miles"`
	Cost_per_gallon float64   `db:"cost_per_gallon"`
	Cost_total      float64   `db:"cost_total"`
	Location        string    `db:"location"`
}

func (db *Data) GetEntryById(id string) (*Entry, error)
	c := Entry
	err := db.Select(&c, "SELECT * FROM log_entry WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *Data) GetEntriesByLogBook(logbook string) (*[]Entry, error) {
	c := []Entry{}
	err := db.Select(&c, "SELECT * FROM log_entry WHERE log_book_id = (SELECT id FROM log_book WHERE id = $1)", logbook)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *Data) InsertEntry(logbook string, miles int64, cost_per_gallon float64, cost_total float64, location string) (*Entry, error) {
	c := Entry
	tx := db.MustBegin()
	result, err := tx.Queryx("INSERT INTO log_entry (log_book_id, miles, cost_per_gallon, cost_total, location) VALUES ((SELECT id FROM log_book WHERE id = $1), $2, $3, $4, $5) RETURNING *", log_book_id, miles, cost_per_gallon, cost_total, location)
	for result.Next() {
		err = result.StructScan(&c)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &c, nil
}
