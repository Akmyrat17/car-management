package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("Connecting to database...")
	time.Sleep(5 * time.Second)

	var err error
	db, err = sql.Open("postgres", connStr) // 🔥 FIXED: no shadowing
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	fmt.Println("Connected to database")
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatal("Error closing database: ", err)
	}
}
