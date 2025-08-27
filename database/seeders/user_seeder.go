package seeders

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/studio-senkou/lentera-cendekia-be/app/models"
)

func SeedUsers(db *sql.DB) error {
	userRepository := models.NewUserRepository(db)

	administrator := &models.User{
		Name:            "Studio Senkou",
		Email:           "studio.senkou@gmail.com",
		Password:        "12345678",
		Role:            "admin",
		EmailVerifiedAt: func() *time.Time { t := time.Now(); return &t }(),
		IsActive:        true,
	}

	if err := userRepository.Create(administrator); err != nil {
		return fmt.Errorf("failed to create administrator user: %w", err)
	}

	return nil
}
