package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (db *Database) Tx() (tx pgx.Tx, cancel context.CancelFunc, err error) {
	ctx := context.Background()
	ctxto, cancel := context.WithTimeout(ctx, 10*time.Second)

	tx, err = db.pool.Begin(ctxto)
	return tx, cancel, err
}
