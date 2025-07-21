package facades

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Database interface {
	GetDB() *sql.DB
	Close() error
}

type DatabaseFacade struct {
	Database *sql.DB
}

func (f *DatabaseFacade) GetDB() *sql.DB {
	return f.Database
}

func (f *DatabaseFacade) Close() error {
	if f.Database != nil {
		return f.Database.Close()
	}
	return nil
}

func (f *DatabaseFacade) Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DATABASE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	f.Database = db
	log.Println("Successfully connected to PostgreSQL database")
	
	return nil
}