package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Workder %d xử lý task %d \n", id, job)
	}
}
func WorkerPool() {
	const numWorker = 3
	const numJobs = 10

	jobs := make(chan int)
	var wg sync.WaitGroup

	// Khởi tạo worker
	for i := 1; i <= numWorker; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Gửi task và channel
	for i := 0; i < numJobs; i++ {
		jobs <- i
	}

	// Báo không còn task mới
	close(jobs)
	wg.Wait()
	fmt.Println("Tất cả task đã hoàn thành")
}

func worker1(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d xử lý task %d\n", id, job)
	}
}
func BAI1() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int)
	var wg sync.WaitGroup

	// Tạo workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker1(i, jobs, &wg)
	}

	// Gửi tasks vào channel
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}

	// Báo rằng không còn task nào được gửi
	close(jobs)

	// Chờ tất cả workers kết thúc
	wg.Wait()

	fmt.Println("Tất cả task đã hoàn thành")
}

func worker2(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range jobs {
		result := num * num

		fmt.Printf(
			"Worker %d: %d^2 = %d\n",
			id,
			num,
			result,
		)
	}
}
func BAI2() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int)
	var wg sync.WaitGroup

	// Khởi chạy workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Gửi các số từ 1 đến 10
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()

	fmt.Println("Đã xử lý xong tất cả các số")
}

func worker3(
	id int,
	jobs <-chan int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d bắt đầu task %d\n", id, job)

		time.Sleep(time.Second)

		fmt.Printf("Worker %d hoàn thành task %d\n", id, job)
	}
}

func BAI3() {
	const numWorkers = 5
	const numJobs = 20

	jobs := make(chan int)
	var wg sync.WaitGroup

	start := time.Now()

	// Tạo 5 workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Gửi 20 tasks
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()

	fmt.Printf(
		"Tất cả task hoàn thành sau: %v\n",
		time.Since(start),
	)

}
func main() {
	BAI3()
}
