package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (db *Database) Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	db.log.Debug("Database.Exec", "sql", sql, "args", args)
	return db.pool.Exec(ctxto, sql, args...)
}

func (db *Database) Query(sql string, args ...any) (pgx.Rows, error) {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	db.log.Debug("Database.Query", "sql", sql, "args", args)
	return db.pool.Query(ctxto, sql, args...)
}

func (db *Database) QueryRow(sql string, args ...any) pgx.Row {
	ctx := context.Background()
	ctxto, _ := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	db.log.Debug("Database.QueryRow", "sql", sql, "args", args)
	return db.pool.QueryRow(ctxto, sql, args...)
}

func (db *Database) ConvertError(err error) error {
	var pgErr *pgconn.PgError

	if err == nil {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("no_rows")
	}

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return errors.New(fmt.Sprintf("unique_violation_on_%s", pgErr.ConstraintName))
		default:
			return err
		}
	} else {
		return err
	}
}
