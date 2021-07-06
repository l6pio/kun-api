package cve

import (
	_ "embed"
	"fmt"
)

//go:embed create-table.sql
var createTableSql string

//go:embed count-all.sql
var countAllSql string

//go:embed insert-report.sql
var insertReportSql string

//go:embed select-all.sql
var selectAllSql string

func CreateTableSql() string {
	return createTableSql
}

func CountAllSql() string {
	return countAllSql
}

func InsertReportSql() string {
	return insertReportSql
}

func SelectAllSql(order string) string {
	return fmt.Sprintf(selectAllSql, order)
}
