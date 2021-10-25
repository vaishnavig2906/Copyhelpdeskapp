package main

import (
	"fmt"
	"log"

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

func InitDB() *sqlx.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user_, password, host, port, dbname)
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

func main() {
	init_router()
}
