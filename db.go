package getresidence

import (
	"database/sql"

	"github.com/efficientgo/core/errors"
)

type grDb struct {
	*sql.DB
}

func newDb(url string) (db grDb, err error) {
	sqlDb, err := sql.Open("libsql", url)
	if err != nil {
		return db, errors.Wrap(err, "open db")
	}

	err = sqlDb.Ping()
	if err != nil {
		return db, errors.Wrap(err, "ping db")
	}

	db = grDb{sqlDb}

	err = db.init()
	if err != nil {
		return db, errors.Wrap(err, "init db")
	}

	return db, nil
}

// sqlite schema
var schema = `
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL
) STRICT
`

func (db grDb) init() error {
	_, err := db.Exec(schema)
	if err != nil {
		return errors.Wrap(err, "exec schema")
	}

	return nil
}

func (db grDb) getOnboarding(id int64) (name, email, phone string, err error) {
	err = db.QueryRow("SELECT name, email, phone FROM sessions WHERE id = :id", sql.Named("id", id)).Scan(&name, &email, &phone)
	if err != nil {
		return "", "", "", errors.Wrap(err, "select session")
	}

	return name, email, phone, nil
}

func (db grDb) newRow() (id int64, err error) {
	err = db.QueryRow(`INSERT INTO sessions (name, email, phone) VALUES ("", "", "") RETURNING id`).Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "insert session")
	}

	return id, nil
}

func (db grDb) setName(id int64, name string) error {
	_, err := db.Exec("UPDATE sessions SET name = :name WHERE id = :id", sql.Named("name", name), sql.Named("id", id))
	if err != nil {
		return errors.Wrap(err, "update name")
	}

	return nil
}

func (db grDb) setEmail(id int64, email string) error {
	_, err := db.Exec("UPDATE sessions SET email = :email WHERE id = :id", sql.Named("email", email), sql.Named("id", id))
	if err != nil {
		return errors.Wrap(err, "update email")
	}

	return nil
}

func (db grDb) setPhone(id int64, phone string) error {
	_, err := db.Exec("UPDATE sessions SET phone = :phone WHERE id = :id", sql.Named("phone", phone), sql.Named("id", id))
	if err != nil {
		return errors.Wrap(err, "update phone")
	}

	return nil
}
