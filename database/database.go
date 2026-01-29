package database

import (
	"database/sql" // SQL database package
	"log"          // Logging package

	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		return nil, err
	}

	// Set connection pool settings if needed
	db.SetMaxOpenConns(25) // set max open connections
	db.SetMaxIdleConns(5)  // set max idle connections

	log.Println("Database connection established")

	return db, nil
}
