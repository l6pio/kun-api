package query

import _ "embed"

//go:embed create-database.sql
var sql string

func CreateDatabaseSQL() string {
	return sql
}
