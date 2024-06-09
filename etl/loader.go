package etl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
	"strings"
)

type loader struct {
	logger *log.CustomLogger
	dbConn *sql.DB
}

// NewLoader creates a new instance of the Loader with the provided database connection.
func NewLoader(logger *log.CustomLogger, dbConn *sql.DB) Loader {
	return &loader{
		logger: logger,
		dbConn: dbConn,
	}
}

// BatchInsert inserts a batch of responses into the PostgreSQL database.
func (l *loader) BatchInsert(ctx context.Context, responses []*model.Response) error {
	// Initialize slices to build the SQL statement
	valueStrings := make([]string, 0, len(responses))     // Slice to hold value placeholders
	valueArgs := make([]interface{}, 0, len(responses)*6) // Slice to hold the actual values

	// Iterate over the responses and construct the values part of the SQL statement
	for i, response := range responses {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, NOW() AT TIME ZONE 'UTC')", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs, response.UserID, response.DeviceType, response.IP, response.DeviceID, response.Locale, response.AppVersion)
	}

	// Join the value strings to form the complete SQL statement
	stmt := fmt.Sprintf("INSERT INTO user_logins (user_id, device_type, masked_ip, masked_device_id, locale, app_version, create_date) VALUES %s",
		strings.Join(valueStrings, ","))

	// Execute the SQL statement with the value arguments
	_, err := l.dbConn.Exec(stmt, valueArgs...)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Failed to execute batch insert with error : %v", err.Error())}
		l.logger.Log(&lm)

		return err
	}

	lm := log.Message{Level: "INFO", Msg: fmt.Sprintf("Successfully inserted a batch to database.")}
	l.logger.Log(&lm)

	return nil
}
