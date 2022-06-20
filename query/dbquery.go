package query

import (
	"fmt"
	"persistence/pgsql"
	"persistence/prop"
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

func (DBQuery) UpdateInto(tableName string, columnNames []string, columnValues []interface{}, conditions []string, conditionsValues []interface{}, conditionals []string, response *interface{}) error {
	var columnQuery = ""
	for index, column := range columnNames {
		if !strings.EqualFold(columnQuery, "") {
			columnQuery = columnQuery + " , "
		}
		columnQuery = columnQuery + fmt.Sprintf("%s %s %s", column, prop.EQUAL, columnValues[index])
	}
	var conditionsQuery = ""
	for index, condition := range conditions {
		if !strings.EqualFold(conditionsQuery, "") {
			conditionsQuery = conditionsQuery + " AND "
		}
		conditionsQuery = conditionsQuery + fmt.Sprintf("%s %s %s", condition, conditionals[index], conditionsValues[index])
	}
	var query = fmt.Sprintf("UPDATE %s SET %s WHERE %s ", tableName, columnQuery, conditionsQuery)
	row := pgsql.ExecuteQueryToOne(query)
	return (*row).Scan(&response)
}
