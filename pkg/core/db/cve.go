package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func FindVulnerabilityById(conf *core.Config, id string) (interface{}, error) {
	session, col, err := GetCol(conf, "vulnerability")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var ret interface{}
	err = col.Find(bson.M{"id": id}).One(&ret)
	return ret, err
}

func FindVulnerabilityByArtifactId(conf *core.Config, id string, page int, order string) (interface{}, error) {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"artId": id}},
		bson.M{"$lookup": bson.M{
			"from":         "vulnerability",
			"localField":   "vulId",
			"foreignField": "id",
			"as":           "vul",
		}},
		bson.M{"$unwind": "$vul"},
		bson.M{"$group": bson.M{"_id": "$vul"}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$_id"}},
	)
	return (&Paging{}).DoPipe(col, stages, page, order)
}

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
