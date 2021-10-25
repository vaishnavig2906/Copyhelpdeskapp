package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user_    = "postgres"
	password = "postgres"
	dbname   = "dpay_helpdesk"
)

func InitDB() (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user_, password, host, port, dbname)
	db, err := sqlx.Open("postgres", dbURL)
	return db, err
}

func main() {
	init_router()
}
