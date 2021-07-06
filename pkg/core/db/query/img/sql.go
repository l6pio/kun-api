package img

import (
	_ "embed"
	"fmt"
)

//go:embed create-table.sql
var createTableSql string

//go:embed count-all.sql
var countAllSql string

//go:embed count-by-id.sql
var countByIdSql string

//go:embed insert-status.sql
var insertStatusSql string

//go:embed pick-id-by-image-id.sql
var pickIdByImageId string

//go:embed select-all.sql
var selectAllSql string

func CreateTableSql() string {
	return createTableSql
}

func CountAllSql() string {
	return countAllSql
}

func CountByIdSql() string {
	return countByIdSql
}

func InsertStatusSql() string {
	return insertStatusSql
}

func PickIdByImageId() string {
	return pickIdByImageId
}

func SelectAllSql(order string) string {
	return fmt.Sprintf(selectAllSql, order)
}
