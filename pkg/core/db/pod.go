package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func SavePod(conf *core.Config, pod *vo.Pod) error {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{"name": pod.Name}, pod)
	return err
}
