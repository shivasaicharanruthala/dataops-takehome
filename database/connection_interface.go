package database

import "database/sql"

type SQLDatabase interface {
	Open() (*sql.DB, error)
}
