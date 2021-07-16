package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func ListAllImages(conf *core.Config, page int, order string) (*Paging, error) {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return (&Paging{}).DoQuery(col.Find(bson.M{}), page, order)
}

func FindImageById(conf *core.Config, id string) (interface{}, error) {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var ret interface{}
	err = col.Find(bson.M{"id": id}).One(&ret)
	return ret, err
}

func FindImageByArtifactId(conf *core.Config, id string, page int, order string) (interface{}, error) {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"artId": id}},
		bson.M{"$lookup": bson.M{
			"from":         "image",
			"localField":   "imgId",
			"foreignField": "id",
			"as":           "img",
		}},
		bson.M{"$unwind": "$img"},
		bson.M{"$group": bson.M{"_id": "$img"}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$_id"}},
	)
	return (&Paging{}).DoPipe(col, stages, page, order)
}

func FindImageByCveId(conf *core.Config, id string, page int, order string) (interface{}, error) {
	session, col, err := GetCol(conf, "cve")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"vulId": id}},
		bson.M{"$lookup": bson.M{
			"from":         "image",
			"localField":   "imgId",
			"foreignField": "id",
			"as":           "img",
		}},
		bson.M{"$unwind": "$img"},
		bson.M{"$group": bson.M{"_id": "$img"}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$_id"}},
	)
	return (&Paging{}).DoPipe(col, stages, page, order)
}

func SaveImage(conf *core.Config, img *vo.Image) error {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{"id": img.Id}, bson.M{"$setOnInsert": img})
	return err
}

func UpdateImageUsage(conf *core.Config, imageId string, status core.ImageEventType) error {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return err
	}
	defer session.Close()

	if status == core.ImageUp {
		return col.Update(bson.M{"id": imageId}, bson.M{"$inc": bson.M{"usage": 1}})
	} else {
		err := col.Update(bson.M{"id": imageId, "usage": bson.M{"$gt": 0}}, bson.M{"$inc": bson.M{"usage": -1}})
		if err != mgo.ErrNotFound {
			return err
		}
		return nil
	}
}
