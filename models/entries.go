package models

import (
	"math/big"
	"time"
)

type Entry struct {
	Id              string    `db:"id"`
	LogBook         string    `db:"log_book_id"`
	Created_at      time.Time `db:"created_at"`
	Miles           *big.Int  `db:"miles"`
	Cost_per_gallon float64   `db:"cost_per_gallon"`
	Cost_total      float64   `db:"cost_total"`
	Location        string    `db:"location"`
}
