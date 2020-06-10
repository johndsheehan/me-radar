package radar

import (
	"bytes"
	"image"
	"image/gif"
	"log"
	"sync"
	"time"

	"github.com/johndsheehan/met-eireann-archive/pkg/queue"
)

// RadarImage interface to fetch radar image
type RadarImage interface {
	Fetch(time.Time) (*image.Paletted, error)
}

// Radar handle rainfall radar fetching, creation, and storage
type Radar struct {
	lock sync.Mutex
	q    *queue.Queue
	gif  []byte

	running bool
	stop    chan struct{}

	RadarImage
}

// NewRadar return new instance of Radar
func NewRadar(history int, fetch RadarImage) *Radar {
	q, err := queue.NewQueue(history)
	if err != nil {
		return nil
	}

	r := Radar{
		sync.Mutex{},
		q,
		nil,
		false,
		make(chan struct{}),
		fetch,
	}

	r.populate()
	return &r
}

// Fetch return latest rainfall radar gif
func (r *Radar) Fetch() []byte {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.gif
}

// Populate Radar object with radar images
func (r *Radar) populate() error {
	history := r.q.MaxLength()

	for i := history; i > 0; i-- {
		d := time.Duration(i * 15)
		then := time.Now().Add(-d * time.Minute)
		gifImg, err := r.RadarImage.Fetch(then)
		if err != nil {
			log.Print(err)
			continue
		}
		r.update(gifImg)
	}

	gifImg, err := r.RadarImage.Fetch(time.Now())
	if err != nil {
		log.Print(err)
		return nil
	}
	r.update(gifImg)

	return nil
}

// Stop watch if it is running
func (r *Radar) Stop() {
	if r.running {
		log.Print("stopping radar")
		r.stop <- struct{}{}
	}

	r.running = false
}

// Watch start a goroutine to periodically check for updated images
func (r *Radar) Watch() {
	if !r.running {
		go r.watch()
	}
}

// update add new gif image to existing gif
func (r *Radar) update(gifImg *image.Paletted) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// store image
	r.q.Push(gifImg)

	// iterate through stored images, create new gif
	g := &gif.GIF{}
	entries, _ := r.q.Entries()
	for _, entry := range entries {
		g.Image = append(g.Image, entry)
		g.Delay = append(g.Delay, 200)
	}

	// final image is shown for longer
	g.Image = append(g.Image, entries[len(entries)-1])
	g.Delay = append(g.Delay, 200)

	var buf bytes.Buffer
	err := gif.EncodeAll(&buf, g)
	if err != nil {
		return err
	}

	r.gif = buf.Bytes()
	return nil
}

func (r *Radar) watch() {
	ticker := time.NewTicker(1 * time.Minute)
	retry := true

	r.running = true

	for {
		select {
		case <-r.stop:
			ticker.Stop()
			break
		case <-ticker.C:
			gifImg, err := r.RadarImage.Fetch(time.Now())
			if err != nil {
				retry = true
				ticker.Stop()

				ticker = time.NewTicker(1 * time.Minute)
				continue
			}

			if retry {
				retry = false
				ticker.Stop()

				ticker = time.NewTicker(15 * time.Minute)
			}

			r.update(gifImg)
		}
	}

	r.running = false
}
