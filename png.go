package main

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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

func textImageCreate(text string, x, y int, bounds image.Rectangle) *image.RGBA {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	img := image.NewRGBA(bounds)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	return img
}

func pngText(img []byte, dateStr, timeStr string) ([]byte, error) {
	decoded, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	dateImg := textImageCreate(dateStr, 10, 30, decoded.Bounds())
	timeImg := textImageCreate(timeStr, 10, 60, decoded.Bounds())

	out := image.NewRGBA(decoded.Bounds())
	draw.Draw(out, decoded.Bounds(), decoded, image.ZP, draw.Src)
	draw.Draw(out, dateImg.Bounds(), dateImg, image.ZP, draw.Over)
	draw.Draw(out, timeImg.Bounds(), timeImg, image.ZP, draw.Over)

	var buf bytes.Buffer
	wtr := bufio.NewWriter(&buf)
	err = png.Encode(wtr, out)
	if err != nil {
		return nil, err
	}
	wtr.Flush()

	return buf.Bytes(), nil
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
