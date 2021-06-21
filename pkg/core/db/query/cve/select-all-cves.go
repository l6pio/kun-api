package cve

import (
	_ "embed"
	"fmt"
)

//go:embed select-all-cves.sql
var selectAllCvesSQL string

func SelectAllCvesSQL(order string) string {
	return fmt.Sprintf(selectAllCvesSQL, order)
}
