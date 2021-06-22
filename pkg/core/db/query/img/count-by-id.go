package img

import (
	_ "embed"
)

//go:embed count-by-id.sql
var countByIdSql string

func CountByIdSql() string {
	return countByIdSql
}
