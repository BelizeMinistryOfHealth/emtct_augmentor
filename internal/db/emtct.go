package db

import (
	"database/sql"
	"fmt"

	"moh.gov.bz/mch/emtct/internal/config"
)

type EmtctDb struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(cnf *config.DbConf) (*EmtctDb, error) {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.Username, cnf.Database, cnf.Password, cnf.Host)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	return &EmtctDb{db}, nil
}
