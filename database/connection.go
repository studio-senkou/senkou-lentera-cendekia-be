package database

import (
	"database/sql"

	"github.com/studio-senkou/lentera-cendekia-be/database/facades"
)

var DB *facades.DatabaseFacade

func InitializeDatabase() error {
	DB = &facades.DatabaseFacade{}
	return DB.Connect()
}

func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB.GetDB()
}
