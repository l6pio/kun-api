package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
)

func ListAllRegistrySettings(conf *core.Config) (*Paging, error) {
	session, col, err := GetCol(conf, "image")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return (&Paging{}).DoQuery(col.Find(bson.M{}), page, order)
}
