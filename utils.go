package main

import (
	_ "image/png"
)

// func timestrings(t time.Time) (string, string, string) {
// 	tsNow := t
// 	minute := (tsNow.Minute() / 15) * 15
// 	tsStr := fmt.Sprintf("%d%02d%02d%02d%02d",
// 		tsNow.Year(),
// 		tsNow.Month(),
// 		tsNow.Day(),
// 		tsNow.Hour(),
// 		minute,
// 	)
// 	dateStr := fmt.Sprintf("%d-%02d-%02d", tsNow.Year(), tsNow.Month(), tsNow.Day())
// 	timeStr := fmt.Sprintf("%02d:%02d", tsNow.Hour(), minute)

// 	return tsStr, dateStr, timeStr
// }

// func snooze(retry bool) bool {
// 	if retry {
// 		time.Sleep(1 * time.Minute)
// 	} else {
// 		time.Sleep(15 * time.Minute)
// 	}
// 	return true
// }

// func fetch(t time.Time) (*image.Paletted, error) {
// 	timestamp, dateStr, timeStr := timestrings(t)
// 	log.Printf("%s, %s, %s", timestamp, dateStr, timeStr)

// 	pngBytes, err := pngFetch(timestamp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pngImg, err := pngText(pngBytes, dateStr, timeStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gifImg, err := pngToGIF(pngImg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return gifImg, nil
// }

// func update(r *Radar) {
// 	retry := true
// 	for {
// 		retry = snooze(retry)
// 		gifImg, err := fetch(time.Now())
// 		if err != nil {
// 			continue
// 		}
// 		r.Update(gifImg)
// 		retry = false
// 	}
// }
