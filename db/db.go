package db

import (
	"database/sql"
	"fmt"

	"github.com/gchaincl/dotsql"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "galaxy"
	dbname = "galaxy_db"
)

type DB struct {
	*sql.DB
}

func InitDB() (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createGalaxyTables(db)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func createGalaxyTables(db *sql.DB) error {
	dot, err := dotsql.LoadFromFile("db/queries.sql")
	if err != nil {
		return err
	}

	_, err = dot.Exec(db, "create-torrents-table")
	if err != nil {
		return err
	}

	_, err = dot.Exec(db, "create-pieces-table")
	if err != nil {
		return err
	}
	return nil
}
