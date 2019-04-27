package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"

	"gocv.io/x/gocv"
)

const (
	urlBase = "http://archive.met.ie/weathermaps/radar2/WEB_radar5_"
)

func pngFetch(timestamp string) ([]byte, error) {
	url := urlBase + timestamp + ".png"
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

func pngText(img []byte, dateStr, timeStr string) ([]byte, error) {
	imgMat, err := gocv.IMDecode(img, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}

	color := color.RGBA{255, 255, 255, 100}
	gocv.PutText(&imgMat, dateStr, image.Point{10, 30}, gocv.FontHersheyComplex, 0.75, color, 2)
	gocv.PutText(&imgMat, timeStr, image.Point{10, 60}, gocv.FontHersheyComplex, 0.75, color, 2)

	return gocv.IMEncode(".png", imgMat)
}

func pngToGIF(pngImg []byte) (*image.Paletted, error) {
	imgData, imgType, err := image.Decode(bytes.NewReader(pngImg))
	if err != nil {
		return nil, err
	}

	if imgType != "png" {
		return nil, errors.New("image type is not png")
	}

	buf := bytes.Buffer{}
	opt := gif.Options{}
	err = gif.Encode(&buf, imgData, &opt)
	if err != nil {
		return nil, err
	}

	img, err := gif.Decode(&buf)
	if err != nil {
		return nil, err
	}

	i, ok := img.(*image.Paletted)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return i, nil
}
