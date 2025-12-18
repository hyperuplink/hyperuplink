package database

import (
	"context"
	"embed"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Database struct {
	log        *slog.Logger
	connection string
	migrations *embed.FS

	poolcfg *pgxpool.Config
	pool    *pgxpool.Pool
}

func New(
	log *slog.Logger,
	connection string,
) (db *Database, err error) {
	db = new(Database)
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

func (db *Database) SetMigrations(migrations *embed.FS) {
	db.migrations = migrations
}

func (db *Database) Startup() (err error) {
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

	sqldb := stdlib.OpenDBFromPool(db.pool)
	psql, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return err
	}

	migrations, err := iofs.New(db.migrations, "migrations")
	if err != nil {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return err
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "postgres", psql)
	if err != nil {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return err
	}

	serr, dberr := m.Close()
	if serr != nil {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return serr
	}
	if dberr != nil {
		psql.Close()
		sqldb.Close()
		db.Shutdown()
		return dberr
	}

	return nil
}

func (db *Database) Shutdown() error {
	db.pool.Close()
	return nil
}
