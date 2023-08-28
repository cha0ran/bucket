package bucket

import (
	"fmt"
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
	b := &bucket{
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

func (b *bucket) start() {
	if b.full {
		for i := 0; i < b.capacity; i++ {
			b.tokensChan <- struct{}{}
		}
	}

	go func() {
		ticker := time.NewTicker(b.duration)
		defer func() {
			ticker.Stop()
			fmt.Println("bucket over.")
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

func (b *bucket) Take() {
	<-b.tokensChan
}

func (b *bucket) Close() {
	b.closeChan <- struct{}{}
	<-b.waitChan
}
