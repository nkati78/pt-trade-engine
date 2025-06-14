package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	dsn string
}

func NewDatabaseConnection(dsn string) *bun.DB {
	return connect(dsn)
}

// const dsn = "postgres://postgres:@localhost:5432/test?sslmode=disable"

// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
// master password RDS QYs1Ecdtv1xvycyo7bGX

func connect(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}
