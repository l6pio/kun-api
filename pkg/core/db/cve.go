package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func SaveCve(conf *core.Config, cve *vo.Cve) error {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{
		"imgId": cve.ImgId,
		"artId": cve.ArtId,
		"vulId": cve.VulId,
	}, bson.M{"$setOnInsert": cve})
	return err
}
