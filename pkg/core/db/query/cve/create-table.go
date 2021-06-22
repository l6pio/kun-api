package cve

import _ "embed"

//go:embed create-table.sql
var createTableSql string

func CreateTableSql() string {
	return createTableSql
}
