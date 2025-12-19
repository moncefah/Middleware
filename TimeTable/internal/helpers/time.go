package helpers

import "time"

// ParseTime parse une date au format RFC3339 (SQLite → Go)
func ParseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
}
