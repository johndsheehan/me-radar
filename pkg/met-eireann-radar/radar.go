package mer

import (
	"fmt"
	"image"
	"time"
)

// RadarFormat type of image
type RadarFormat int8

const (
	// ARCHIVE old style image
	ARCHIVE RadarFormat = 0
	// CURRENT recent style image
	CURRENT RadarFormat = 1
)

type meRadar interface {
	fetch(string) ([]byte, error)
}

// MERadar met eireann radar
type MERadar struct {
	format RadarFormat
	meRadar
}

// NewMERadar return new instance
func NewMERadar(rf RadarFormat) (*MERadar, error) {
	var me meRadar

	me = newArchive()
	if rf == CURRENT {
		me = newCurrent()
	}

	return &MERadar{
		format:  rf,
		meRadar: me,
	}, nil
}

// Fetch gif format radar image
func (m MERadar) Fetch(t time.Time) (*image.Paletted, error) {
	timestamp, dateStr, timeStr := m.timestrings(t)

	pngBytes, err := m.fetch(timestamp)
	if err != nil {
		return nil, err
	}

	pngImg, err := pngText(pngBytes, dateStr, timeStr)
	if err != nil {
		return nil, err
	}

	gifImg, err := pngToGIF(pngImg)
	if err != nil {
		return nil, err
	}

	return gifImg, nil
}

func (m MERadar) timestrings(t time.Time) (string, string, string) {
	tsNow := t
	minute := (tsNow.Minute() / 15) * 15

	tsStrFormat := "%d%02d%02d%02d%02d"
	if m.format == CURRENT {
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
