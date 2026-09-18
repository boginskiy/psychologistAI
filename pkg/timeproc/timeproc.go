package timeproc

import (
	"net/http"
	"time"
)

// ConvertTimeUtcToLocal. Конвертация время UTC -> Local Time
func ConvertTimeUtcToLocal(timeUTC *time.Time, timezone string) string {
	// Format:
	// "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation(timezone)
	return timeUTC.In(loc).Format(time.RFC3339)
}

// ConvertTimeUtcToLocalByRequest.
func ConvertTimeUtcToLocalByRequest(timeUTC *time.Time, r *http.Request) *time.Time {
	// Take Tz
	tz := r.Header.Get("Accept-Timezone")
	if tz == "" {
		return timeUTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return timeUTC
	}
	newTime := timeUTC.In(loc)
	return &newTime
}
