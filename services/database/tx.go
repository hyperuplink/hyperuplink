package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Tx struct {
	Tx        pgx.Tx
	Ctx       context.Context
	CtxTO     context.Context
	CtxCancel context.CancelFunc
}

func (db *Database) Tx() (tx *Tx, err error) {
	tx = new(Tx)
	tx.Ctx = context.Background()
	tx.CtxTO, tx.CtxCancel = context.WithTimeout(tx.Ctx, 10*time.Second)

	if tx.Tx, err = db.pool.Begin(tx.CtxTO); err != nil {
		// tx.CtxCancel()
		return nil, err
	}

	return tx, nil
}

func (tx *Tx) Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	return tx.Tx.Exec(ctxto, sql, args...)
}

func (tx *Tx) Query(sql string, args ...any) (pgx.Rows, error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	return tx.Tx.Query(ctxto, sql, args...)
}

func (tx *Tx) QueryRow(sql string, args ...any) pgx.Row {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	return tx.Tx.QueryRow(ctxto, sql, args...)
}

func (tx *Tx) Commit() (err error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	return tx.Tx.Commit(ctxto)
}

func (tx *Tx) End() (err error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	err = tx.Tx.Rollback(ctxto)
	// tx.CtxCancel()
	return err
}
