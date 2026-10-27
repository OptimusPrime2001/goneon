package main

import (
	"fmt"
	"sync"
	"testing"
)

type Product struct {
	mu    sync.RWMutex
	price int
}

func (p *Product) getPrice() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.price
}
func (p *Product) setPrice() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.price *= 2
	return p.price
}

func TestRwmutext(t *testing.T) {
	product := Product{
		price: 100,
	}

	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Go(func() {
			fmt.Println("Get price:", product.getPrice())
		})
	}
	wg.Go(
		func() {
			product.setPrice()
		})
	wg.Wait()
}
