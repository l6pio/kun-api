package cve

import _ "embed"

//go:embed insert-report.sql
var insertReportSql string

func InsertReportSql() string {
	return insertReportSql
}
