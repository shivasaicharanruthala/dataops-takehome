package store

import (
	"database/sql"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
)

type login struct {
	dbConn *sql.DB
}

func New(dbConn *sql.DB) Login {
	return &login{
		dbConn: dbConn,
	}
}

func (l login) Get(filter map[string]string) ([]model.Response, error) {

	return nil, nil
}
