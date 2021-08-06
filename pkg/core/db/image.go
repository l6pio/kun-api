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

func ImageExists(conf *core.Config, id string) (bool, error) {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return false, err
	}
	defer session.Close()

	count, err := col.Find(bson.M{"id": id}).Count()
	return count > 0, err
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

func FindImageTimelineById(conf *core.Config, id string) ([]interface{}, error) {
	session, col, err := GetCol(conf, "image-timeline")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var ret []interface{}
	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"imageId": id}},
		bson.M{"$group": bson.M{
			"_id": bson.M{
				"$subtract": []bson.M{
					{"$toLong": "$timestamp"},
					{"$mod": []interface{}{bson.M{"$toLong": "$timestamp"}, 60 * 1000}},
				},
			},
			"count": bson.M{"$sum": "$status"},
		}},
		bson.M{"$addFields": bson.M{"timestamp": "$_id"}},
	)
	err = col.Pipe(stages).All(&ret)
	return ret, err
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

func SaveImageTimelineEvent(conf *core.Config, timeline *vo.ImageTimeline) error {
	session, col, err := GetCol(conf, "image-timeline")
	if err != nil {
		return err
	}
	defer session.Close()

	return col.Insert(timeline)
}

func UpdateImagePods(conf *core.Config, event *vo.ImageTimeline) error {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return err
	}
	defer session.Close()

	if event.Status == vo.ImageUp {
		return col.Update(bson.M{"id": event.ImageId}, bson.M{"$inc": bson.M{"pods": 1}})
	} else {
		err := col.Update(bson.M{"id": event.ImageId, "pods": bson.M{"$gt": 0}}, bson.M{"$inc": bson.M{"pods": -1}})
		if err != mgo.ErrNotFound {
			return err
		}
		return nil
	}
}

func CleanImagePods(conf *core.Config) error {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.UpdateAll(bson.M{}, bson.M{"$set": bson.M{"pods": 0}})
	return err
}
