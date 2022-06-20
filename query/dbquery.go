package query

import (
	"fmt"
	"persistence/pgsql"
	"strings"
)

type DBQuery struct{}

func (DBQuery) ExecuteQuery(query string, response *interface{}) error {
	row := pgsql.ExecuteQueryToOne(query)
	return (*row).Scan(&response)
}
func (DBQuery) ExistsByColumn(tableName string, column string, data string, response *bool) error {
	var query = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = '%s')", tableName, column, data)
	row := pgsql.ExecuteQueryToOne(query)
	return (*row).Scan(&response)
}

func (DBQuery) FindAllBy(tableName string, columnName string, value string, response *interface{}) error {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", tableName, columnName, value)
	row, err := pgsql.ExecuteQueryToMany(query)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func (DBQuery) InsertInto(tableName string, columnNames []string, values []interface{}, response *interface{}) error {
	var query = fmt.Sprintf("INSERT INTO %s (%s) VALUES", tableName, strings.Join(columnNames, ","))
	var queryValues = query + " (%d) RETURNING *"
	row, err := pgsql.ExecuteQueryParams(queryValues, values)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}
