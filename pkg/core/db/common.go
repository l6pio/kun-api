package db

import (
	_ "embed"
	"gopkg.in/mgo.v2"
	"l6p.io/kun/api/pkg/core"
	"time"
)

func GetCol(conf *core.Config, name string) (*mgo.Session, *mgo.Collection, error) {
	session, err := mgo.DialWithInfo(dialInfo(conf.MongoAddr, conf.DatabaseName, conf.MongoUser, conf.MongoPass))
	if err != nil {
		return nil, nil, err
	}

	col := session.DB(conf.DatabaseName).C(name)
	return session, col, nil
}

func dialInfo(dbAddr string, dbName string, username string, password string) *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    []string{dbAddr},
		Direct:   false,
		Timeout:  time.Second * 15,
		Database: dbName,
		Source:   "admin",
		Username: username,
		Password: password,
	}
}
