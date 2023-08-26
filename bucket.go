package bucket

import (
	"time"
)

type bucket struct {
	capacity   int
	quantum    int
	duration   time.Duration
	tokensChan chan struct{}
	closeChan  chan struct{}
	waitChan   chan struct{}
	full       bool
}

func NewBucket(duration time.Duration, capacity int, quantum int, full bool) *bucket {
	return &bucket{
		capacity:   capacity,
		quantum:    quantum,
		duration:   duration,
		tokensChan: make(chan struct{}, capacity),
		closeChan:  make(chan struct{}, 1),
		waitChan:   make(chan struct{}, 1),
		full:       full,
	}
}

func (b *bucket) Start() {
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

func (b *bucket) Stop() {
	b.closeChan <- struct{}{}
}

func (b *bucket) Take() {
	<-b.tokensChan
}

func (b *bucket) Wait() {
	<-b.waitChan
}
