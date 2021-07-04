package img

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/query/img"
	"time"
)

const (
	StatusUp   = 1
	StatusDown = 0
)

func Exists(conn *sql.DB, imageId string) (bool, error) {
	stmt, err := conn.Prepare(img.CountByIdSql())
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(imageId)
	if err != nil {
		return false, err
	}

	rows.Next()

	var count int64
	if err := rows.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func Status(conn *sql.DB, imageId string, image string, status int) (string, error) {
	id, err := db.RunTx(conn, img.InsertStatusSql(), func(stmt *sql.Stmt) (interface{}, error) {
		id := uuid.New().String()
		_, err := stmt.Exec(id, imageId, image, status, time.Now())
		return id, err
	})
	if err != nil {
		return "", err
	}
	log.Infof("update image '%s' status to UP", imageId)
	return id.(string), nil
}

func PickId(conn *sql.DB, imageId string) (string, error) {
	stmt, err := conn.Prepare(img.PickIdByImageId())
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(imageId, time.Now().Truncate(time.Hour))
	if err != nil {
		return "", err
	}

	rows.Next()

	var id string
	if err := rows.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
