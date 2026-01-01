package datetime

import (
	"time"

	"github.com/studio-senkou/lentera-cendekia-be/app/models"
)

func ParseDateOnly(dateStr string) (models.DateOnly, error) {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return models.DateOnly{}, err
	}
	return models.DateOnly(parsedTime), nil
}

func ParseTimeOnly(timeStr string) (models.TimeOnly, error) {
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		parsedTime, err = time.Parse("15:04", timeStr)
		if err != nil {
			return models.TimeOnly{}, err
		}
	}
	return models.TimeOnly(parsedTime), nil
}
