```go
package main

import (
	"fmt"
	"github.com/cha0ran/bucket"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	b := bucket.NewBucket(time.Minute, 60, 1, false)
	defer b.Close()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(taskId int) {
			defer wg.Done()
			b.Take()
			fmt.Printf("Task number: %d", taskId)
		}(i)
	}

	wg.Wait()
}
```