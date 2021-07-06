package img

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve/vo/api"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/db/query/img"
	"time"
)

const (
	StatusUp   = 1
	StatusDown = 0
)

func List(conf *core.Config, page int, order string) (*db.Paging, error) {
	ret, err := (&db.Paging{
		Page: page,
		DoCount: func() (*sql.Rows, error) {
			return conf.DbConn.Query(img.CountAllSql())
		},
		DoQuery: func(from int, size int) (*sql.Rows, error) {
			return conf.DbConn.Query(img.SelectAllSql(order), from, size)
		},
		Convert: func(rows *sql.Rows) []interface{} {
			ret := make([]interface{}, 0)
			for rows.Next() {
				var id string
				var name string
				var size int64
				var artCount int64
				var vulCount int64

				if err := rows.Scan(&id, &name, &size, &artCount, &vulCount); err != nil {
					log.Error(err)
				}

				ret = append(ret, api.Image{
					Id:       id,
					Name:     name,
					Size:     size,
					ArtCount: artCount,
					VulCount: vulCount,
				})
			}
			return ret
		},
	}).Do()
	if err != nil {
		return nil, err
	}
	return ret, err
}

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
