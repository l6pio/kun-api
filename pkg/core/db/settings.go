package db

import (
	"gopkg.in/mgo.v2/bson"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func ListAllRegistryAuthSettings(conf *core.Config) (ret []vo.RegistryAuth, err error) {
	session, col, err := GetCol(conf, "setting")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	err = col.Find(bson.M{"type": vo.RegistryAuthType}).All(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SaveRegistryAuthSettings(conf *core.Config, auths []vo.RegistryAuth) error {
	session, col, err := GetCol(conf, "setting")
	if err != nil {
		return err
	}
	defer session.Close()

	for _, auth := range auths {
		_, err = col.Upsert(bson.M{"authority": auth.Authority}, auth)
		if err != nil {
			return err
		}
	}
	return nil
}
