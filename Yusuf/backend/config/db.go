package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("db_host"),
		os.Getenv("db_user"),
		os.Getenv("db_password"),
		os.Getenv("db_name"),
		os.Getenv("db_port"),
	)
	log.Println("DSN:", dsn) // DEBUG LOG

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("DB unreachable: %v", err)
	}
}
