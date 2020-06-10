package mer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
)

type constantVars struct {
	archiveURL string

	apiURL string
	host   string

	headers map[string]string

	xRange []int
	yRange []int

	tileOffsets map[string]tileOffset

	currentBasePNG image.Image
}

func setConstantVars() constantVars {
	archiveURL := "http://archive.met.ie/weathermaps/radar2/WEB_radar5_"

	apiURL := "https://api.met.ie/api/maps/radar"
	host := "api.met.ie"

	headers := make(map[string]string)
	headers["Origin"] = "https://www.met.ie"
	headers["Referer"] = "https://www.met.ie"
	headers["Connection"] = "keep-alive"
	headers["Cache-Control"] = "no-cache"
	headers["Pragma"] = "no-cache"

	xRange := []int{59, 60, 61, 62, 63}
	yRange := []int{40, 41, 42, 43}

	offsets := make(map[string]tileOffset)
	for i, x := range xRange {
		for j, y := range yRange {
			n := fmt.Sprintf("%d%d", x, y)
			offsets[n] = tileOffset{i * 256, j * 256}
		}
	}

	basePNG, err := base64.StdEncoding.DecodeString(b64CurrentBaseImg)
	if err != nil {
		log.Fatal(err)
	}

	rdr := bytes.NewReader(basePNG)
	imagePNG, _, err := image.Decode(rdr)
	if err != nil {
		log.Fatal(err)
	}
	return constantVars{
		archiveURL:     archiveURL,
		apiURL:         apiURL,
		host:           host,
		headers:        headers,
		xRange:         xRange,
		yRange:         yRange,
		tileOffsets:    offsets,
		currentBasePNG: imagePNG,
	}
}

var (
	constants = setConstantVars()
)
