package mer

import (
	"fmt"
	"time"
)

func timestrings(t time.Time) (string, string, string) {
	tsNow := t
	minute := (tsNow.Minute() / 15) * 15
	tsStr := fmt.Sprintf("%d%02d%02d%02d%02d",
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
