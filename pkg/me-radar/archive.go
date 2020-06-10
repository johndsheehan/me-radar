package mer

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type meArchive struct {
}

func newArchive() *meArchive {
	return &meArchive{}
}

func (m meArchive) fetch(timestamp string) ([]byte, error) {
	url := constants.archiveURL + timestamp + ".png"
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
