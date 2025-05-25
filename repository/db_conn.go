package repository

import (
	"2do.com/config"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

func ConnectDB() *sql.DB {
	// Get the configuration from the config package
	dbConfig := config.LoadDBConfig()

	// Create connection string using the configuration
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName,
	)

	// Open database connection
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to reach the database: %v", err)
	}

	log.Println("Database connection successfully established.")
	return db
}
