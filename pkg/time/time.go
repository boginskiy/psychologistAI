package time

import "time"

// Конвертация время UTC -> Local Time
func ConvertTimeUtcToLocal(tm time.Time, timezone string) string {
	// Format:
	// "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation(timezone)
	return tm.In(loc).Format(time.RFC3339)
}
