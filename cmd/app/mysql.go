package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eveisesi/athena/internal/mysqldb"
)

func makeMySQL(cfg config) (*sql.DB, error) {

	m := cfg.MySQL

	db, err := mysqldb.New(context.TODO(), m.Host, m.Port, m.User, m.Pass, m.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql db: %w", err)
	}

	return db, nil
}
