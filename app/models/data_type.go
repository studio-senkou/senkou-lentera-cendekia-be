package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type TimeOnly time.Time

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	timeStr := time.Time(t).Format("15:04:05")
	return []byte(`"` + timeStr + `"`), nil
}

func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsedTime, err := time.Parse("15:04:05", str)
	if err != nil {
		parsedTime, err = time.Parse("15:04", str)
		if err != nil {
			return err
		}
	}
	*t = TimeOnly(parsedTime)
	return nil
}

func (t *TimeOnly) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*t = TimeOnly(v)
		return nil
	case []byte:
		parsedTime, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = TimeOnly(parsedTime)
		return nil
	case string:
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*t = TimeOnly(parsedTime)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", value)
	}
}

func (t TimeOnly) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05"), nil
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	dateStr := time.Time(d).Format("2006-01-02")
	return []byte(`"` + dateStr + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = DateOnly(parsedTime)
	return nil
}

func (d *DateOnly) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*d = DateOnly(parsedTime)
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*d = DateOnly(parsedTime)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot scan %T into DateOnly", value)
	}
}

func (d DateOnly) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}
