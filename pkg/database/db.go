package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"gitlab.com/hieuhani/permitbox/internal/config"
	"gitlab.com/hieuhani/permitbox/pkg/atomicity"
	"gitlab.com/hieuhani/permitbox/pkg/shutdown"
	"io/fs"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type (
	GetDbFunc func(ctx context.Context) bun.IDB
)

func New(cfg config.DbConfig, tasks *shutdown.Tasks, migrationSource fs.FS) (GetDbFunc, *atomicity.DbAtomicExecutor, error) {
	emptyAtomicExecutor := &atomicity.DbAtomicExecutor{}
	completeDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)

	conn, err := sql.Open("postgres", completeDsn)
	if err != nil {
		return nil, emptyAtomicExecutor, err
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(2 * time.Hour)

	db := bun.NewDB(conn, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if cfg.SqlDebugEnabled {
		db.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(true),
				bundebug.WithVerbose(true),
			),
		)
	}
	if err := conn.Ping(); err != nil {
		return nil, emptyAtomicExecutor, err
	}

	err = MigrationUp(cfg.DbName, conn, migrationSource)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		break
	case err != nil:
		return nil, emptyAtomicExecutor, err
	}

	getDbFunc := func(ctx context.Context) bun.IDB {
		if tx := atomicity.ContextGetTx(ctx); tx.Tx != nil {
			return tx
		}
		return db
	}

	tasks.AddShutdownTask(
		func(_ context.Context) error {
			return db.Close()
		},
	)

	return getDbFunc, &atomicity.DbAtomicExecutor{DB: db}, nil
}

func MigrationUp(dbName string, db *sql.DB, migrations fs.FS) error {
	iofsDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", iofsDriver, dbName, driver)
	if err != nil {
		return err
	}

	return migrator.Up()
}
