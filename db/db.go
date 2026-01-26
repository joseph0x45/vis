package db

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/sad"
)

type Conn struct {
	db *sqlx.DB
}

func (c *Conn) Close() {
	c.db.Close()
}

func Connect(opts sad.DBConnectionOptions) *Conn {
	if opts.Reset {
		log.Println("Starting with fresh database")
	}
	db, err := sad.OpenDBConnection(opts, migrations)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database at", opts.DatabasePath)
	return &Conn{db}
}

func rollbackTx(tx *sqlx.Tx, originalErr error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Join(originalErr, err)
	}
	return originalErr
}
