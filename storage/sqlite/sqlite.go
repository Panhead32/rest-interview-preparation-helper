package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func InitDatabase(pathName string) (*Storage, error) {
	if pathName == "" {
		log.Fatal("Database path is required")
	}

	db, err := sql.Open("sqlite3", pathName)

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt, err := os.ReadFile("./storage/db_schema.sql")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlStmt))

	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) DB() *sql.DB {
	if s.db == nil {
		log.Fatal("Database not initialized")
	}

	return s.db
}
