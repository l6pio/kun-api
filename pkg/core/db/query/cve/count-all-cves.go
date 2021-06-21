package cve

import _ "embed"

//go:embed count-all-cves.sql
var countAllCvesSQL string

func CountAllCvesSQL() string {
	return countAllCvesSQL
}
