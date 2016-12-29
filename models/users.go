package models

import (
	"time"

	gouuid "github.com/satori/go.uuid"
)

type Uuid gouuid.UUID

type User struct {
	Id             string    `db:"id"`
	Uname          string    `db:"uname"`
	Role           string    `db:"role"`
	Digest         []byte    `db:"digest"`
	Email          string    `db:"email"`
	Last_online_at time.Time `db:"last_online_at"`
	Created_at     time.Time `db:"created_at"`
}

//func (db *Data) queryOne(query string) (*User, error)

//func (db *Data) queryMany(query string) (*[]User, error)

func (db *Data) InsertUser(name string, digest []byte, email string) (*User, error) {
	var u User
	tx := db.MustBegin()
	result, err := tx.Queryx("INSERT INTO users (uname, digest, email) VALUES ($1, $2, $3) RETURNING *", name, string(digest[:]), email)
	for result.Next() {
		err = result.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (db *Data) GetUserByEmail(email string) (*User, error) {
	var u User
	tx := db.MustBegin()
	result, err := tx.Queryx("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		err = result.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &u, nil
}

func (db *Data) LoginUser(email string, digest []byte) (*User, error) {
	var u User
	tx := db.MustBegin()
	result, err := tx.Queryx("SELECT * FROM users WHERE email = $1 AND digest = $2", email, digest)
	for result.Next() {
		err = result.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &u, nil
}

func (db *Data) GetUserById(id string) (*User, error) {
	var u User
	tx := db.MustBegin()
	result, err := tx.Queryx("SELECT * FROM users WHERE id = $1", id)
	for result.Next() {
		err = result.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &u, nil
}
