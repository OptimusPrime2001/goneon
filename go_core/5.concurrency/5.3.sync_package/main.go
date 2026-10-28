package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += 1
}
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}
func BAI1() {
	var wg sync.WaitGroup
	counter := Counter{
		count: 0,
	}
	for i := 0; i < 100; i++ {

		wg.Go(func() {
			for i := 0; i < 1000; i++ {
				counter.Increment()

			}
		})
	}
	wg.Wait()
	fmt.Println("Count", counter.Value())
}

func BAI2() {
	numbers := make([]int, 0, 10)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		numbers = append(numbers, i)
	}
	results := make([]int, len(numbers))
	for i, value := range numbers {
		wg.Go(func() {
			results[i] = value * value
		})
	}

	wg.Wait()
	fmt.Println(results)
}

type Store struct {
	rwmu sync.RWMutex
	data map[string]int
}

func (s *Store) Set(key string, value int) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()
	s.data[key] = value
}
func (s *Store) Get(key string) (int, bool) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	value, ok := s.data[key]
	return value, ok

}
func BAI3() {
	store := Store{
		data: map[string]int{},
	}
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Go(func() {
			store.Set("Key "+strconv.Itoa(i), i)
		})
	}
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			value, ok := store.Get("Key " + strconv.Itoa(i))
			if ok {
				fmt.Println("Đọc dữ liệu thành công", value)
			} else {
				fmt.Println("Không có dữ liệu cho key này")
			}
		})
	}
	wg.Wait()
	fmt.Println("Main finished")
}

type BankAccount struct {
	mu      sync.Mutex
	balance int
}

func (ba *BankAccount) Deposit(amount int) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.balance += amount
}
func (ba *BankAccount) Withdraw(amount int) bool {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	if ba.balance >= amount {
		ba.balance -= amount
		return true
	}
	return false

}
func BAI4() {
	bankAccount := BankAccount{
		balance: 10000,
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Go(func() {
			bankAccount.Deposit(100)
			bankAccount.Withdraw(50)
		})
	}
	wg.Wait()
	fmt.Println("Số dư cuối cùng", bankAccount.balance)
}

func main() {
	BAI1()
	BAI2()
	BAI3()
	BAI4()
}
