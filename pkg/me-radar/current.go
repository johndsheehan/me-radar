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
	"sync"
)

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

type tile struct {
	name string
	xpos int
	ypos int
	data []byte
}

type tileOffset struct {
	x int
	y int
}

type meCurrent struct {
	icon image.Image
}

func newCurrent() *meCurrent {
	mec := &meCurrent{
		icon: image.NewRGBA(image.Rect(0, 0, 1, 1)),
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", constants.iconURL, nil)
	if err != nil {
		log.Print(err)
		return mec
	}

	rsp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return mec
	}
	defer rsp.Body.Close()

	iconData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Print(err)
		return mec
	}

	icon, err := png.Decode(bytes.NewReader(iconData))
	if err != nil {
		log.Print(err)
		return mec
	}

	mec.icon = icon
	return mec
}

func (m meCurrent) fetch(timestamp string) ([]byte, error) {
	baseURL, endpoint, err := fetchRadarURL(timestamp)
	if err != nil {
		return nil, err
	}

	tileImg, err := fetchTileImage(baseURL, endpoint)
	if err != nil {
		return nil, nil
	}

	basePNG := constants.currentBasePNG

	out := image.NewRGBA(basePNG.Bounds())
	draw.Draw(out, basePNG.Bounds(), basePNG, image.ZP, draw.Over)
	draw.Draw(out, basePNG.Bounds(), tileImg, image.Point{X: 80, Y: 0}, draw.Over)
	draw.Draw(out, basePNG.Bounds(), m.icon, image.ZP, draw.Over)

	var buf bytes.Buffer
	wtr := bufio.NewWriter(&buf)
	err = png.Encode(wtr, out)
	if err != nil {
		return nil, err
	}
	wtr.Flush()

	return buf.Bytes(), nil
}

func fetchRadarURL(timestamp string) (string, string, error) {
	urlList, err := fetchURLList()
	if err != nil {
		return "", "", err
	}

	var entries []radarInfoEntry
	json.Unmarshal(urlList, &entries)

	url, end := "", ""
	for _, e := range entries {
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

func fetchTile(url string, xpos, ypos int, ch chan tile, wg *sync.WaitGroup) {
	defer wg.Done()

	name := fmt.Sprintf("%d%d", xpos, ypos)
	t := tile{
		name: name,
		xpos: xpos,
		ypos: ypos,
		data: nil,
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		ch <- t
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		ch <- t
		return
	}

	t.data = make([]byte, len(data))
	copy(t.data, data)

	ch <- t
}

func fetchTileImage(baseURL, endpoint string) (image.Image, error) {
	// https://devcdn.metweb.ie/api/maps/radar/202005301315/61/41/7/1590844879
	total := len(constants.xRange) * len(constants.yRange)
	tileChan := make(chan tile, total)
	var wg sync.WaitGroup

	for _, x := range constants.xRange {
		for _, y := range constants.yRange {
			url := fmt.Sprintf("%s/%d/%d/7/%s", baseURL, x, y, endpoint)

			wg.Add(1)
			go fetchTile(url, x, y, tileChan, &wg)
		}
	}

	wg.Wait()
	close(tileChan)

	tileImg := image.NewRGBA(image.Rect(0, 0, (256 * 5), (256 * 4)))
	for t := range tileChan {
		decoded, err := png.Decode(bytes.NewReader(t.data))
		if err != nil {
			log.Print(err)
			continue
		}

		o := constants.tileOffsets[t.name]
		bbox := image.Rect(o.x, o.y, o.x+256, o.y+256)
		draw.Draw(tileImg, bbox.Bounds(), decoded, image.ZP, draw.Over)
	}

	return tileImg, nil
}

func fetchURLList() ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", constants.apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Host = constants.host
	for k, v := range constants.headers {
		req.Header.Add(k, v)
	}

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
