package facades

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Database interface {
	GetDB() *sql.DB
	Close() error
	Transaction(fn func(tx *sql.Tx) error) error
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

func (f *DatabaseFacade) Transaction(fn func(tx *sql.Tx) error) error {
	tx, err := f.Database.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is shadowed in return, but we handle it via named return or manual check
		}
	}()

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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