package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second):
			fmt.Println("Worker stopped after", time.Second)
		case <-ctx.Done():
			fmt.Println("Worker stopped: ", ctx.Err())
			switch {
			case errors.Is(ctx.Err(), context.Canceled):
				fmt.Println("Chủ động cancel")
				if err := context.Cause(ctx); err != nil {
					fmt.Println("Dừng do reason cụ thể là:", err)
				}
				// bị hủy chủ động
			case errors.Is(ctx.Err(), context.DeadlineExceeded):
				fmt.Println("Dừng do deadline time")
			}
			return
		default:
			if dl, ok := ctx.Deadline(); ok {
				fmt.Println("Còn lại:", time.Until(dl))
			}
			fmt.Println("Worker is working...")
			time.Sleep(time.Millisecond * 500)
		}
	}
}
func fetchData(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Data fetched")
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// Các loại context
func contextCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx)
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(time.Millisecond * 500)
	fmt.Println("Main finished")
}
func contextTimeout() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := fetchData(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Success")
}
func contextDeadline() {
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	go worker(ctx)
	defer cancel()
	time.Sleep(6 * time.Second)
	fmt.Println("Main finished")
}
func contextCancelCause() {
	ctx, cancel := context.WithCancelCause(context.Background())
	go worker(ctx)
	time.Sleep(5 * time.Second)
	cancel(errors.New("Database mất kết nối"))
	time.Sleep(time.Millisecond * 500)
	fmt.Println("Main finished")
}
func main() {
	// BAI1()
	// BAI2()
	// BAI3()
	BAI4()

}
