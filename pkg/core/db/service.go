package db

import (
	"database/sql"
	_ "embed"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core/db/query"
	"l6p.io/kun/api/pkg/core/db/query/cve"
	"l6p.io/kun/api/pkg/core/db/query/img"
)

func Init(conn *sql.DB) {
	_, err := conn.Exec(query.CreateDatabaseSQL())
	if err != nil {
		log.Fatalf("Create database error: %v", err)
	}

	_, err = conn.Exec(cve.CreateTableSql())
	if err != nil {
		log.Fatalf("Create 'cve' table error: %v", err)
	}

	_, err = conn.Exec(img.CreateTableSql())
	if err != nil {
		log.Fatalf("Create 'img' table error: %v", err)
	}
}

func Connect() *sql.DB {
	//TODO: put this into config
	conn, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	Ping(conn)
	log.Info("Database is connected")
	return conn
}

func Ping(conn *sql.DB) {
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Fatalf("Database error '%d': %s\n%s",
				exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Fatalf("Database error: %v", err)
		}
	}
}

func RunTx(conn *sql.DB, query string, action func(stmt *sql.Stmt) (interface{}, error)) (interface{}, error) {
	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ret, err := action(stmt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ret, nil
}
