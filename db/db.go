package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresSQLStorage(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
