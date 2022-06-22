package repo

import (
	"fmt"
	"gitlab.falabella.com/falabella-retail/txd/boss/libraries/golang/persistence/client"
	"strings"
)

func GetPgRepo(cnt *client.PgClient) DBQuery {
	return &pgRepo{cnt}
}

type pgRepo struct {
	cnt *client.PgClient
}

func (r *pgRepo) ExecuteQuery(query string, response *interface{}) error {
	row := r.cnt.RunQueryRow(query)
	return (*row).Scan(&response)
}

func (r *pgRepo) FindAllBy(tableName string, condition DBCondition, response *interface{}) error {
	op, err := parseOperator(condition.Operator)
	if err != nil {
		return err
	}
	var query = fmt.Sprintf("SELECT * FROM %s WHERE %s %s '%s'", tableName, condition.FieldName, op, condition.Value)
	row, err := r.cnt.RunQueryRows(query)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func (r *pgRepo) FindByConditions(tableName string, conditions []DBCondition, response *interface{}) error {
	var conditionsQuery = ""
	for _, condition := range conditions {
		if !strings.EqualFold(conditionsQuery, "") {
			conditionsQuery = conditionsQuery + " AND "
		}
		op, err := parseOperator(condition.Operator)
		if err != nil {
			return err
		}
		conditionsQuery = conditionsQuery + fmt.Sprintf("(%s %s '%s')", condition.FieldName, op, condition.Value)
	}
	var query = fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, conditionsQuery)
	row, err := r.cnt.RunQueryRows(query)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func (r *pgRepo) ExistsBy(tableName string, condition DBCondition, response *bool) error {
	op, err := parseOperator(condition.Operator)
	if err != nil {
		return err
	}
	var query = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s %s '%s')", tableName, condition.FieldName, op, condition.Value)
	row, err := r.cnt.RunQueryRows(query)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func (r *pgRepo) InsertInto(tableName string, columnNames []string, values []interface{}, response *interface{}) error {
	var query = fmt.Sprintf("INSERT INTO %s (%s) VALUES", tableName, strings.Join(columnNames, ","))
	var queryValues = query + " (%d) RETURNING *"
	row, err := r.cnt.RunQueryArgs(queryValues, values)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func (r *pgRepo) Update(tableName string, columnValues []DBValue, conditions []DBCondition, response *interface{}) error {
	var columnQuery = ""
	for _, column := range columnValues {
		if !strings.EqualFold(columnQuery, "") {
			columnQuery = columnQuery + " , "
		}
		columnQuery = columnQuery + fmt.Sprintf("%s = %s", column.FieldName, column.Value)
	}
	var conditionsQuery = ""
	for _, condition := range conditions {
		if !strings.EqualFold(conditionsQuery, "") {
			conditionsQuery = conditionsQuery + " AND "
		}
		op, err := parseOperator(condition.Operator)
		if err != nil {
			return err
		}
		conditionsQuery = conditionsQuery + fmt.Sprintf("%s %s %s", condition.FieldName, op, condition.Value)
	}
	var query = fmt.Sprintf("UPDATE %s SET %s WHERE %s ", tableName, columnQuery, conditionsQuery)
	row, err := r.cnt.RunQueryRows(query)
	if err != nil {
		return err
	}
	return (*row).Scan(&response)
}

func parseOperator(op DBOperator) (string, error) {
	switch op {
	case EQUAL:
		return "=", nil
	case NOTEQUAL:
		return "<>", nil
	case GREATER:
		return ">", nil
	case LESS:
		return "<", nil
	case GREATERE:
		return ">=", nil
	case LESSE:
		return "<=", nil
	case LIKE:
		return "LIKE", nil
	case NOTLIKE:
		return "NOT LIKE", nil
	case IN:
		return "IN", nil
	case NOTIN:
		return "NOT IN", nil
	case ISNULL:
		return "IS NULL", nil
	case ISNOTNULL:
		return "IS NOT NULL", nil
	default:
		return "", fmt.Errorf("operator '%s' not suported", op)
	}
}
