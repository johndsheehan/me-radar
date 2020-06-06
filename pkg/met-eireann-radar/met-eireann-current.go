package mer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type MECurrentConfig struct {
}

type MECurrent struct {
}

func NewMECurrent(cfg *MECurrentConfig) (*MECurrent, error) {
	return &MECurrent{}, nil
}

func (m MECurrent) Fetch(t time.Time) (*image.Paletted, error) {
	timestamp, dateStr, timeStr := timestrings(t, current)
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

type tile struct {
	name string
	xpos int
	ypos int
	data []byte
}

func (m MECurrent) fetch(timestamp string) ([]byte, error) {
	baseURL, endpoint, err := fetchRadarURL(timestamp)
	if err != nil {
		return nil, err
	}

	baseFile, err := os.Open("base.png")
	if err != nil {
		log.Fatal(err)
	}
	defer baseFile.Close()

	basePNG, _, err := image.Decode(baseFile)
	if err != nil {
		log.Fatal(err)
	}

	// https://devcdn.metweb.ie/api/maps/radar/202005301315/61/41/7/1590844879
	xRange := []int{59, 60, 61, 62, 63}
	yRange := []int{40, 41, 42, 43}
	total := len(xRange) * len(yRange)

	tileChan := make(chan tile, total)
	var wg sync.WaitGroup

	for _, x := range xRange {
		for _, y := range yRange {
			url := fmt.Sprintf("%s/%d/%d/7/%s", baseURL, x, y, endpoint)

			wg.Add(1)
			go fetchTile(url, x, y, tileChan, &wg)
		}
	}

	wg.Wait()
	close(tileChan)

	type offset struct {
		x int
		y int
	}
	moffset := make(map[string]offset)
	for i, x := range xRange {
		for j, y := range yRange {
			n := fmt.Sprintf("%d%d", x, y)
			moffset[n] = offset{i * 256, j * 256}
		}
	}

	out := image.NewRGBA(basePNG.Bounds())
	draw.Draw(out, basePNG.Bounds(), basePNG, image.ZP, draw.Over)

	for t := range tileChan {
		decoded, err := png.Decode(bytes.NewReader(t.data))
		if err != nil {
			log.Print(err)
			continue
		}

		o := moffset[t.name]
		bbox := image.Rect(o.x, o.y, o.x+256, o.y+256)
		draw.Draw(out, bbox.Bounds(), decoded, image.ZP, draw.Over)
	}

	var buf bytes.Buffer
	wtr := bufio.NewWriter(&buf)
	err = png.Encode(wtr, out)
	if err != nil {
		return nil, err
	}
	wtr.Flush()

	return buf.Bytes(), nil
}

func fetchTile(url string, xpos, ypos int, ch chan tile, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("fetching %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}

	name := fmt.Sprintf("%d%d", xpos, ypos)

	t := tile{
		name: name,
		xpos: xpos,
		ypos: ypos,
		data: make([]byte, len(data)),
	}
	copy(t.data, data)

	ch <- t
}

type radarInfoEntry struct {
	Src              string `json:"src"`
	DateAndTime      string `json:"dateAndTime"`
	DayName          string `json:"dayName"`
	MapDate          string `json:"mapDate"`
	MapTime          string `json:"mapTime"`
	FullMapTimestamp string `json:"fullMapTimestamp"`
	ToolTipDate      string `json:"toolTipDate"`
	ModifiedTime     int    `json:"modifiedTime"`
	Server           string `json:"server"`
	DayIndex         int    `json:"dayIndex"`
}

func fetchRadarURL(timestamp string) (string, string, error) {
	urlList, err := fetchURLList()
	if err != nil {
		return "", "", err
	}

	var entries []radarInfoEntry
	json.Unmarshal(urlList, &entries)
	log.Print(entries)

	url := ""
	end := ""
	for _, e := range entries {
		log.Printf("checking %s against %s\n", timestamp, e.DateAndTime)
		if timestamp == e.DateAndTime {
			url = fmt.Sprintf("%s/api/maps/radar/%s", e.Server, e.Src)
			end = fmt.Sprintf("%d", e.ModifiedTime)
			break
		}
	}

	if url == "" {
		return "", "", errors.New("url not found")
	}

	return url, end, nil
}

func fetchURLList() ([]byte, error) {
	client := &http.Client{}

	url := "https://devcdn.metweb.ie/api/maps/radar"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Host = "api.met.ie"
	req.Header.Add("Origin", "https://www.met.ie")
	req.Header.Add("Referer", "https://www.met.ie")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	urlList, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return urlList, nil
}
