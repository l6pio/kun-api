package cve

import (
	_ "embed"
	"fmt"
)

//go:embed select-all.sql
var selectAllSql string

func SelectAllSql(order string) string {
	return fmt.Sprintf(selectAllSql, order)
}
