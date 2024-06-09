package database

import (
	"database/sql"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"os"
)

type dbConn struct {
	logger *log.CustomLogger
}

// New returns a new instance of SQLDatabase interface by creating and returning a pointer to dbConn struct.
func New(logger *log.CustomLogger) SQLDatabase {
	return &dbConn{
		logger: logger,
	}
}

// Open initializes and opens a connection to the database using environment variables for configuration.
func (dbo *dbConn) Open() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbDriver := os.Getenv("DRIVER_NAME")

	// Create the connection string using the retrieved environment variables.
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize DB connection
	db, err := sql.Open(dbDriver, connectionStr)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Database initilization failed with error %v", err.Error())}
		dbo.logger.Log(&lm)

		return nil, err
	}

	return db, nil
}
