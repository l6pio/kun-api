package img

import (
	_ "embed"
)

//go:embed pick-id-by-image-id.sql
var pickIdByImageId string

func PickIdByImageId() string {
	return pickIdByImageId
}
