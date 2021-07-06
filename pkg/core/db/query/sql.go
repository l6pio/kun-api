package query

import _ "embed"

//go:embed create-database.sql
var createDatabaseSql string

func CreateDatabaseSQL() string {
	return createDatabaseSql
}
