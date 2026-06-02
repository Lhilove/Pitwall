package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() {
	connStr := os.Getenv("db")

	if connStr == "" {
		log.Fatal("Environment variable 'db' is not set")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	DB = db
	fmt.Println("Connected to PostgreSQL")
}
