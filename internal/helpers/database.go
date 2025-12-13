package helpers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "timetable.db?_foreign_keys=on")

	if err != nil {
		return nil, err
	}

	// IMPORTANT : une seule connexion pour SQLite
	db.SetMaxOpenConns(1)

	return db, nil
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		logrus.Errorf("error closing db : %s", err.Error())
	}
}
