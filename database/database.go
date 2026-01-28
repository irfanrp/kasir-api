package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(ConnectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", ConnectionString)
	if err != nil {
		return nil, err
	}

	// Test Connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}
