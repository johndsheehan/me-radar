package mer

import (
	"fmt"
	"time"
)

type radarFormat int8

const (
	archive radarFormat = 0
	current radarFormat = 1
)

func timestrings(t time.Time, format radarFormat) (string, string, string) {
	tsNow := t
	minute := (tsNow.Minute() / 15) * 15

	tsStrFormat := "%d%02d%02d%02d%02d"
	if format == current {
		tsStrFormat = "%d-%02d-%02d %02d:%02d"
	}

	tsStr := fmt.Sprintf(tsStrFormat,
		tsNow.Year(),
		tsNow.Month(),
		tsNow.Day(),
		tsNow.Hour(),
		minute,
	)

	dateStr := fmt.Sprintf("%d-%02d-%02d", tsNow.Year(), tsNow.Month(), tsNow.Day())
	timeStr := fmt.Sprintf("%02d:%02d", tsNow.Hour(), minute)

	return tsStr, dateStr, timeStr
}
