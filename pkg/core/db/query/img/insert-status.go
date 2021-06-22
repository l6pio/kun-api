package img

import _ "embed"

//go:embed insert-status.sql
var insertStatusSql string

func InsertStatusSql() string {
	return insertStatusSql
}
