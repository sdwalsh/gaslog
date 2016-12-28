package models

type Log struct {
	Id   string `db:"id"`
	User string `db:"user_id"`
}
