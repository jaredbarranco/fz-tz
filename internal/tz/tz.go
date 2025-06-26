package tz

import (
	"fmt"
	"time"
)

type TzData struct {
	DatabaseTimezone string `json:"databaseTimezone"`
	AbbreviatedTimezone string `json:"abbreviatedTimezoneCode"`
	UtcOffsetHours string `json:"utcOffsetHours"`
	UtcOffsetSeconds int `json:"utcOffsetSeconds"`
	LocalTime string `json:"localTime"`
	Iso8601 string `json:"iso8601"`
}

func GetTzOffset(iso8601Str string, tzStr string) TzData {
	// Example input time
	input := iso8601Str // Try changing to a summer date like "2025-06-15T15:04:05"
	t, err := time.Parse("2006-01-02T15:04:05", input)
	if err != nil {
		panic(err)
	}

	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		panic(err)
	}

	localTime := t.In(loc)
	tzAbbrev, offset := localTime.Zone() // returns name (e.g. "PST") and offset in seconds

	// Format offset like -0700 or -0800
	offsetStr := fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60)

	return TzData{
		DatabaseTimezone: tzStr,
		AbbreviatedTimezone: tzAbbrev,
		UtcOffsetHours: offsetStr,
		UtcOffsetSeconds: offset,
		LocalTime: localTime.String(),
		Iso8601: localTime.Format(time.RFC3339),
	}
}
