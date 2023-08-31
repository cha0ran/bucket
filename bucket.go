package bucket

import (
	"time"
)

type Bucket struct {
	capacity   int
	quantum    int
	duration   time.Duration
	tokensChan chan struct{}
	closeChan  chan struct{}
	waitChan   chan struct{}
	full       bool
}

func NewBucket(duration time.Duration, capacity int, quantum int, full bool) *Bucket {
	b := &Bucket{
		capacity:   capacity,
		quantum:    quantum,
		duration:   duration,
		tokensChan: make(chan struct{}, capacity),
		closeChan:  make(chan struct{}, 1),
		waitChan:   make(chan struct{}, 1),
		full:       full,
	}

	b.start()

	return b
}

func (b *Bucket) start() {
	if b.full {
		for i := 0; i < b.capacity; i++ {
			b.tokensChan <- struct{}{}
		}
	}

	go func() {
		ticker := time.NewTicker(b.duration)
		defer func() {
			ticker.Stop()
			b.waitChan <- struct{}{}
		}()

	Loop:
		for {
			select {
			case <-ticker.C:
				for i := 0; i < b.quantum; i++ {
					select {
					case b.tokensChan <- struct{}{}:
					default:
						break
					}
				}
			case <-b.closeChan:
				break Loop
			}
		}
	}()
}

func (b *Bucket) Take() {
	<-b.tokensChan
}

func (b *Bucket) Close() {
	b.closeChan <- struct{}{}
	<-b.waitChan
}
