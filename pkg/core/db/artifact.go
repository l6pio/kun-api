package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func FindArtifactByImageId(conf *core.Config, id string, page int, order string) (interface{}, error) {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"imgId": id}},
		bson.M{"$lookup": bson.M{
			"from":         "artifact",
			"localField":   "artId",
			"foreignField": "id",
			"as":           "art",
		}},
		bson.M{"$unwind": "$art"},
		bson.M{"$group": bson.M{"_id": "$art"}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$_id"}},
	)
	return (&Paging{}).DoPipe(col, stages, page, order)
}

func FindArtifactByCveId(conf *core.Config, id string, page int, order string) (interface{}, error) {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"vulId": id}},
		bson.M{"$lookup": bson.M{
			"from":         "artifact",
			"localField":   "artId",
			"foreignField": "id",
			"as":           "art",
		}},
		bson.M{"$unwind": "$art"},
		bson.M{"$group": bson.M{"_id": "$art"}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$_id"}},
	)
	return (&Paging{}).DoPipe(col, stages, page, order)
}

func SaveArtifact(conf *core.Config, art *vo.Artifact) error {
	session, col, err := GetCol(conf, "artifact")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{"id": art.Id}, bson.M{"$setOnInsert": art})
	return err
}
