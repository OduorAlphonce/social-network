package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func nullableUUID(value sql.NullString) (*uuid.UUID, error) {
	if !value.Valid || value.String == "" {
		return nil, nil
	}
	parsed, err := uuid.FromString(value.String)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func nullableTime(value sql.NullString) (*time.Time, error) {
	if !value.Valid || value.String == "" {
		return nil, nil
	}
	parsed, err := parseSQLiteTime(value.String)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, format := range formats {
		parsed, err := time.Parse(format, value)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported sqlite time %q", value)
}
