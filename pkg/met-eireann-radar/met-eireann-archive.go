package mer

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MEArchiveConfig struct {
	URLBase string
}

type MEArchive struct {
	urlBase string
}

func NewMEArchive(cfg *MEArchiveConfig) (*MEArchive, error) {
	urlBase := "http://archive.met.ie/weathermaps/radar2/WEB_radar5_"
	if cfg.URLBase != "" {
		urlBase = cfg.URLBase
	}

	return &MEArchive{
		urlBase: urlBase,
	}, nil
}

func (m MEArchive) Fetch(t time.Time) (*image.Paletted, error) {
	timestamp, dateStr, timeStr := timestrings(t)
	log.Printf("%s, %s, %s", timestamp, dateStr, timeStr)

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

func (m MEArchive) fetch(timestamp string) ([]byte, error) {
	url := m.urlBase + timestamp + ".png"
	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("image not found")
	}

	pngBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return pngBytes, nil
}

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
