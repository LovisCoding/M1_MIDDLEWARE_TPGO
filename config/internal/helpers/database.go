package helpers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var NatsConn *nats.Conn

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:config.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		logrus.Errorf("error closing db : %s", err.Error())
	}
}
