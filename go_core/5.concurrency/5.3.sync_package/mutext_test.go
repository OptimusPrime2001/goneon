package main

import (
	"sync"
	"testing"
)

var buss int

func increment(n int, mu *sync.Mutex) {
	for i := 0; i < n; i++ {
		// Lấy khóa, nếu khóa đang được goroutine khác giữ thì phải chờ goroutine đó nhả khóa qua mu.Unlock()
		mu.Lock()
		buss++
		// Sau khi set xong, thì nhả khóa để goroutine khác có thể lấy và set dữ liệu
		mu.Unlock()
	}
}
func TestCounterRace(t *testing.T) {
	buss = 0

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	go func() {
		defer wg.Done()
		increment(10000, &mu)
	}()

	go func() {
		defer wg.Done()
		increment(10000, &mu)
	}()

	wg.Wait()

	if buss != 20000 {
		t.Errorf("Expected 20000, got %d", buss)
	}
}
