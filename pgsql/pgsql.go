package pgsql

import (
	"github.com/jackc/pgx/v4"
	"log"
)

func ExecuteQueryToMany(query string) (*pgx.Rows, error) {
	log.Printf("Executing query: %s\n", query)
	return getPgCnx().runQueryRows(query)
}

func ExecuteQueryToOne(query string) *pgx.Row {
	log.Printf("Executing query: %s\n", query)
	return getPgCnx().runQueryRow(query)
}

func ExecuteQueryParams(query string, params ...interface{}) (*pgx.Rows, error) {
	log.Printf("Executing query: %s, with params: %v\n", query, params)
	return getPgCnx().runQueryArgs(query, params)
}
