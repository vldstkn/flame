package main

import (
	"database/sql"
	"flag"
	"flame/internal/config"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}

func main() {
	var db string
	flag.StringVar(&db, "db", "", "database name (account/swipe)")
	flag.Parse()

	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}
	conf := config.LoadConfig("configs", mode)

	switch db {
	case "swipe":
		dbConn, err := initDB(conf.Database.Swipes.Dsn)
		if err != nil {
			log.Fatalf("Failed to connect to swipes db: %v", err)
		}
		defer dbConn.Close()

		if err := goose.Up(dbConn, fmt.Sprintf("./migrations/%s", db)); err != nil {
			log.Fatalf("Error migrating %s: %v", db, err)
		}
	case "account":
		dbConn, err := initDB(conf.Database.Account.Dsn)
		if err != nil {
			log.Fatalf("Failed to connect to account db: %v", err)
		}
		defer dbConn.Close()

		if err := goose.Up(dbConn, fmt.Sprintf("./migrations/%s", db)); err != nil {
			log.Fatalf("Error migrating %s: %v", db, err)
		}
	default:
		swipeDB, err := initDB(conf.Database.Swipes.Dsn)
		if err != nil {
			log.Fatalf("Failed to connect to swipes db: %v", err)
		}
		defer swipeDB.Close()

		accountDB, err := initDB(conf.Database.Account.Dsn)
		if err != nil {
			log.Fatalf("Failed to connect to account db: %v", err)
		}
		defer accountDB.Close()

		if err := goose.Up(swipeDB, "./migrations/swipe"); err != nil {
			log.Fatalf("Error migrating swipe: %v", err)
		}
		if err := goose.Up(accountDB, "./migrations/account"); err != nil {
			log.Fatalf("Error migrating account: %v", err)
		}
	}
}
