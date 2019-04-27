package main

import (
	"bytes"
	"image"
	"image/gif"
	"sync"
)

// Radar handle rainfall radar fetching, creation, and storage
type Radar struct {
	lock sync.Mutex
	q    *Queue
	gif  []byte
}

// NewRadar return new instance of Radar
func NewRadar(history int) *Radar {
	q, err := NewQueue(history)
	if err != nil {
		return nil
	}

	return &Radar{
		sync.Mutex{},
		q,
		nil,
	}
}

// Fetch return latest rainfall radar gif
func (r *Radar) Fetch() []byte {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.gif
}

// Update add new gif image to existing gif
func (r *Radar) Update(gifImg *image.Paletted) error {
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
