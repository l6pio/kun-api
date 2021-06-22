package cve

import _ "embed"

//go:embed count-all.sql
var countAllSql string

func CountAllSql() string {
	return countAllSql
}
