package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func BAI1() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Go(func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Working stopped")
				return
			case <-ticker.C:
				fmt.Println("Working...")
			}
		}
	})
	time.Sleep(time.Second * 2)
	cancel()
	wg.Wait()

}
func doWork(ctx context.Context) error {
	fmt.Println("Working start")
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Work done")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func BAI2() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	err := doWork(ctx)

	if err != nil {
		fmt.Println("Task failed:", err)
		return
	}

}
func BAI3() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int, ctx context.Context) {
			defer wg.Done()
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					fmt.Printf("Worker %d is Working....\n", id)
				case <-ctx.Done():
					fmt.Printf("Worker %d is stopped\n", id)
					return
				}
			}
		}(i, ctx)
	}
	time.Sleep(time.Second * 2)
	fmt.Println("Canceling workers...")
	cancel()
	wg.Wait()
	fmt.Println("All workers stopped")
}
func BAI4() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	jobs := make(chan int)
	results := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						fmt.Printf("Worker %d finished\n", id)
						return
					}
					fmt.Printf("Worker %d handling job %d\n", id, job)
					// Mô phỏng task cần 1 giây.
					select {
					case <-time.After(time.Second):
						results <- job
						fmt.Printf("Worker %d completed job %d\n", id, job)

					case <-ctx.Done():
						fmt.Printf("Worker %d canceled job %d\n", id, job)
						return
					}
				case <-ctx.Done():
					fmt.Println("Worker stop receive task")
					return
				}
			}

		}(i)
	}

	// Gửi tasks; dừng gửi nếu context hết hạn.
	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped sending jobs:", ctx.Err())
			close(jobs)
			wg.Wait()
			fmt.Printf("Completed tasks: %d\n", len(results))
			return

		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()
	fmt.Printf("Completed tasks: %d\n", len(results))

}
