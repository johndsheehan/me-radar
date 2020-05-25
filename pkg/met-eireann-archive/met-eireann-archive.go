package mea

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
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

func textImageCreate(text string, x, y int, bounds image.Rectangle) *image.RGBA {
	txtImg := image.NewRGBA(bounds)

	draw.Draw(txtImg, bounds, image.Transparent, image.ZP, draw.Src)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Println("truetype.Parse failed")
		return nil
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(24)
	ctx.SetClip(bounds)
	ctx.SetDst(txtImg)
	ctx.SetSrc(image.Black)

	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	ctx.DrawString(text, point)

	return txtImg
}

func pngText(img []byte, dateStr, timeStr string) ([]byte, error) {
	decoded, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	dateImg := textImageCreate(dateStr, 90, 30, decoded.Bounds())
	timeImg := textImageCreate(timeStr, 90, 60, decoded.Bounds())

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
