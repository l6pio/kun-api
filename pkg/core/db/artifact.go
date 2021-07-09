package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func SaveArtifact(conf *core.Config, art *vo.Artifact) error {
	session, col, err := GetCol(conf, "artifact")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{"name": art.Name, "version": art.Version}, bson.M{"$setOnInsert": art})
	return err
}
