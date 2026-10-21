package main

import (
	"fmt"
	"sync"
	"time"
)

// WaitGroupDemo — "chờ TẤT CẢ goroutine con xong".
//   wg.Add(n)   → báo trước có n việc cần chờ (PHẢI chạy TRƯỚC `go`)
//   wg.Done()   → báo 1 việc đã xong (thường defer trong goroutine)
//   wg.Wait()   → main đứng chờ tới khi counter về 0
func WaitGroupDemo() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // Add TRƯỚC khi go — sai chỗ này là ăn race/panic
		go func(id int) {
			defer wg.Done() // đảm bảo chạy kể cả khi có panic được recover
			fmt.Printf("  worker %d đang làm...\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	wg.Wait() // block tới khi cả 3 Done()
	fmt.Println("  tất cả worker đã xong")
}

// WorkerPool — 3 worker luân phiên lấy job từ kênh `jobs`,
// trả kết quả vào kênh `results`.
//
//	jobs ──► [worker x3] ──► results
//
// Số worker CỐ ĐỊNH → không tạo 1000 goroutine cho 1000 job.
func WorkerPool() {
	const (
		numWorkers = 3
		numJobs    = 6
	)

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	// 1. Khởi động ĐÚNG numWorkers goroutine (chúng chạy lại nhiều lần)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs { // lặp tới khi jobs bị close
				fmt.Printf("    worker %d xử lý job %d\n", id, j)
				result := j * j
				time.Sleep(30 * time.Millisecond) // giả lập công việc nặng
				results <- result
			}
		}(w)
	}

	// 2. Gửi toàn bộ job (buffered → không block lâu)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // báo "hết job" → worker thoát vòng for

	// 3. Chờ worker xong, RỒI mới đóng results
	wg.Wait()
	close(results)

	// 4. Thu kết quả
	total := 0
	for r := range results {
		total += r
	}
	fmt.Printf("  tổng kết quả = %d (%d workers, %d jobs)\n", total, numWorkers, numJobs)
}

// WorkerPoolDynamic — job đến ĐỘNG (không biết trước số lượng).
// Thêm job bất kỳ lúc nào; main chỉ cần đóng jobs khi THẬT SỰ hết.
func WorkerPoolDynamic() {
	jobs := make(chan int)
	results := make(chan int)

	var wg sync.WaitGroup
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				results <- j * 10
			}
		}(w)
	}

	// producer: tự quyết định lúc nào gửi, lúc nào đóng
	go func() {
		for _, v := range []int{5, 8, 13, 21} {
			jobs <- v
		}
		close(jobs)
	}()

	// Đóng results trong goroutine RIÊNG: wg.Wait() không phải main
	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println("  tổng kết quả (dynamic) =", sum) // 5+8+13+21 = 47 → 470
}
