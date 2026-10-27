package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var count atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			for j := 0; j < 100000; j++ {
				count.Add(1)
			}
		})
	}

	wg.Wait()
	fmt.Println("Count:", count.Load())
}
