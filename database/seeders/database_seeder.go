package main

import (
	"database/sql"
	"fmt"

	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/database"
)

func main() {
	if err := database.InitializeDatabase(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer database.CloseDatabase()
	
	db := database.GetDB()
	
	if err := Seed(db); err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Database seeded successfully")
}

func Seed(db *sql.DB) error {
	userRepository := models.NewUserRepository(db)
	
	administrator := &models.User{
		Name:            "Studio Senkou",
		Email:           "studio.senkou@gmail.com",
		Password:        "12345678",
		Role:            "admin",
		IsActive:        true,
	}

	if err := userRepository.Create(administrator); err != nil {
		return fmt.Errorf("failed to create administrator user: %w", err)
	}
	
	return nil
}
