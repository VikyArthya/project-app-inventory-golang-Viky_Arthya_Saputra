package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var err error
	connStr := "user=postgres password=superadmin dbname=inventory sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func Close() {
	DB.Close()
}
