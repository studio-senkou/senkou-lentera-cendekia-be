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

	administrator2 := &models.User{
		Name:            "Lentera Cendekia",
		Email:           "lbblenteracendekia@gmail.com",
		Password:        "12345678",
		Role:            "admin",
		EmailVerifiedAt: func() *time.Time { t := time.Now(); return &t }(),
		IsActive:        true,
	}

	if err := userRepository.Create(administrator2); err != nil {
		return fmt.Errorf("failed to create second administrator user: %w", err)
	}

	return nil
}
