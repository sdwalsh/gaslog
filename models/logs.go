package models

type Log struct {
	Id   string `db:"id"`
	User string `db:"user_id"`
}

func (db *Data) GetLogById(id string) (*Entry, error)
	c := Entry
	err := db.Select(&c, "SELECT * FROM log_book WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *Data) GetLogsByUser(logbook string) (*[]Entry, error) {
	c := []Entry{}
	err := db.Select(&c, "SELECT * FROM log_book WHERE user = (SELECT id FROM user WHERE id = $1)", logbook)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *Data) InsertLog(logbook string, miles int64, cost_per_gallon float64, cost_total float64, location string) (*Entry, error) {
	c := Entry
	tx := db.MustBegin()
	result, err := tx.Queryx("INSERT INTO log_book (user_id) VALUES (SELECT id FROM user WHERE id = $1) RETURNING *")
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
