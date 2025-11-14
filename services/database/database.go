package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type IDatabase interface {
	Startup() error
	Shutdown() error
}

type Database struct {
	log        *slog.Logger
	connection string

	poolcfg *pgxpool.Config
	pool    *pgxpool.Pool
}

func New(log *slog.Logger, connection string) (*Database, error) {
	var err error

	db := new(Database)
	db.log = log
	db.connection = connection

	if db.poolcfg, err = pgxpool.ParseConfig(connection); err != nil {
		return nil, err
	}

	db.poolcfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxUUID.Register(c.TypeMap())
		return nil
	}

	return db, nil
}

func (db *Database) Startup() error {
	var err error

	if db.pool, err = pgxpool.NewWithConfig(
		context.Background(),
		db.poolcfg,
	); err != nil {
		return err
	}

	var greeting string
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.pool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		db.Shutdown()
		return err
	}

	return nil
}

func (db *Database) Shutdown() error {
	db.pool.Close()
	return nil
}
