package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Fan-in/ Fan-out
func BAI1() {
	var wg sync.WaitGroup

	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for num := range jobs {
				result := num * num
				fmt.Printf("Worker %d processed %d\n", workerID, num)
				results <- result
			}
		}(w)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)

	for value := range results {
		fmt.Println("Result:", value)
	}
}
func stage1() <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		for i := 1; i <= 20; i++ {
			result <- i
		}
	}()

	return result
}
func stage2(in <-chan int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for value := range in {
			result <- value * 3
		}
	}()
	return result
}
func stage3(in <-chan int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for value := range in {
			if value%2 == 0 {
				result <- value
			}
		}
	}()
	return result
}

// Pipline
func BAI2() {
	stag1 := stage1()
	stag2 := stage2(stag1)
	stag3 := stage3(stag2)
	for n := range stag3 {
		fmt.Println(n)
	}
}

func BAI3() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	jobs := make(chan int, 20)
	results := make(chan int, 20)

	const numWorkers = 4

	var wg sync.WaitGroup

	// Fan-out: khởi chạy 4 workers.
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d stopped: %v\n",
						workerID, ctx.Err())
					return

				case job, ok := <-jobs:
					if !ok {
						return
					}

					fmt.Printf("Worker %d handling job %d\n",
						workerID, job)

					// Mô phỏng xử lý job trong 300ms.
					select {
					case <-time.After(300 * time.Millisecond):
						results <- job
						fmt.Printf("Worker %d completed job %d\n",
							workerID, job)

					case <-ctx.Done():
						fmt.Printf("Worker %d canceled job %d\n",
							workerID, job)
						return
					}
				}
			}
		}(i)
	}

	// Gửi 20 jobs.
	for i := 1; i <= 20; i++ {
		jobs <- i
	}
	close(jobs)

	// Chờ tất cả workers kết thúc.
	wg.Wait()
	close(results)

	// Fan-in: đếm những job đã hoàn thành.
	completed := 0
	for range results {
		completed++
	}

	fmt.Printf("Completed tasks: %d\n", completed)
}
func main() {
	BAI3()
}
