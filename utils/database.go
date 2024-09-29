package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(dsn string) error {
	var err error
	// Connect to database
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("Failed to open a DB connection: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}

	fmt.Println("Database connection established!")
	return nil
}
