package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type DB interface {
	RunTransaction(ctx context.Context, opts *sql.TxOptions, f func(tx *sql.Tx) error) error
	Close() error
}

type db sql.DB

func New(driver, dsn string) (DB, error) {
	openedDatabase, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return (*db)(openedDatabase), nil
}

func (database *db) RunTransaction(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	// Cancel the context (rolling back the tx if it was not committed) on exit.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := (*sql.DB)(database).BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (database *db) Close() error {
	return (*sql.DB)(database).Close()
}
