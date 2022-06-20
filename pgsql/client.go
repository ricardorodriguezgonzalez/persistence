package pgsql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/url"
	"persistence/prop"
	"sync"
	"time"
)

const ulrConnectionString = "postgresql://%s:%s@%s:%s/%s?sslmode=%s"

var (
	instance *PgCnx
	once     sync.Once
)

type PgCnx struct {
	ctx     context.Context
	poolCnx *pgxpool.Pool
}

func (db *PgCnx) runQueryRows(query string) (*pgx.Rows, error) {
	rows, err := db.poolCnx.Query(db.ctx, query)
	if err != nil {
		log.Printf("Error running query (%s): %v\n", query, err)
	}
	return &rows, err
}
func (db *PgCnx) runQueryArgs(query string, args ...interface{}) (*pgx.Rows, error) {
	rows, err := db.poolCnx.Query(db.ctx, query, args)
	if err != nil {
		log.Printf("Error running query (%s): %v\n", query, err)
	}
	return &rows, err
}

func (db *PgCnx) runQueryRow(query string) *pgx.Row {
	row := db.poolCnx.QueryRow(db.ctx, query)
	return &row
}

func getPgCnx() *PgCnx {
	once.Do(func() {
		ctx := context.Background()
		poolCnx := createPoolConnection(ctx)
		instance = &PgCnx{ctx, poolCnx}
	})
	return instance
}

func createPoolConnection(ctx context.Context) *pgxpool.Pool {
	pgProp := prop.GetDbProp()
	urlCnx := fmt.Sprintf(
		ulrConnectionString, url.QueryEscape(pgProp.DbUser), url.QueryEscape(pgProp.DbPassword),
		url.QueryEscape(pgProp.DbHost), url.QueryEscape(pgProp.DbPort), url.QueryEscape(pgProp.DbName),
		url.QueryEscape(pgProp.DbSslMode),
	)

	cnxCfg, err := pgxpool.ParseConfig(urlCnx)
	if err != nil {
		log.Fatalf("Error creating pool config: %v\n", err)
	}

	cnxCfg.MaxConns = pgProp.DbMaxConns
	cnxCfg.AfterRelease = func(cnx *pgx.Conn) bool {
		time.Sleep(pgProp.DbWaitAfterQuery)
		return true
	}

	dbPool, err := pgxpool.ConnectConfig(ctx, cnxCfg)
	if err != nil {
		log.Fatalf("Error creating connection: %v\n", err)
	}
	return dbPool
}
