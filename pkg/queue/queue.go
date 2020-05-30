package queue

import (
	"errors"
	"image"
	"sync"
)

// Queue hold the previous radar images
type Queue struct {
	lock      sync.Mutex
	entries   []*image.Paletted
	maxLength int
}

// NewQueue return new queue instance for radar images
func NewQueue(maxLength int) (*Queue, error) {
	if maxLength < 1 {
		return nil, errors.New("queue size < 1")
	}

	return &Queue{
		sync.Mutex{},
		make([]*image.Paletted, 0),
		maxLength,
	}, nil
}

// Entries return all stored radar images
func (q *Queue) Entries() ([]*image.Paletted, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.entries, nil
}

// MaxLength return maximum capacity of queue
func (q *Queue) MaxLength() int {
	return q.maxLength
}

// Push add newest radar image
func (q *Queue) Push(entry *image.Paletted) {
	q.lock.Lock()
	defer q.lock.Unlock()

	l := len(q.entries)
	if l == q.maxLength {
		q.entries = q.entries[1:l]
	}

	q.entries = append(q.entries, entry)
}
