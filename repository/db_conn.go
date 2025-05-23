package repository

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=taskpassword dbname=task sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to reach the database: %v", err)
	}
	log.Println("Database connection successfully established.")
	return db
}
