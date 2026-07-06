package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"xn--gckvb8fzb.com/hyperuplink/errs"
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
		return errs.ErrNoRows
	}

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("%w_%s", errs.ErrUniqueViolationOn, pgErr.ConstraintName)
		default:
			return err
		}
	} else {
		return err
	}
}
