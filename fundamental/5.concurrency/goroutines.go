package main

import (
	"fmt"
	"sync"
	"time"
)

func sayHello() {
	fmt.Println("Say hello from goroutine")
}

var counter int

func increment() {
	counter++
}

func calculate(ch chan int) {
	result := 10 + 20
	ch <- result
}

// Goroutines — entry point (được main.go gọi).
// Go KHÔNG chờ goroutine nào: main() thoát là TOÀN BỘ goroutine bị dừng.
// Ở đây phải time.Sleep thì sayHello mới kịp in ra — và đây là cách SAI.
// Cách đúng: dùng channel hoặc sync.WaitGroup (xem Channels / Worker Pools).
func Goroutines() {
	go sayHello()
	// Chờ 1s để sayHello chạy xong ở goroutine trước khi hàm kết thúc
	time.Sleep(time.Second)

	ch := make(chan int)
	go calculate(ch)
	result := <-ch
	fmt.Println("Result from channel:", result)

	go increment()
	go increment()

	ConcurrencyVsParallelism()

	GoroutineLeak()
	DataRace()
	LoopCapture()
	fmt.Println("Hello from main")
}

// ConcurrencyVsParallelism:
//
//	CONCURRENT  = nhiều việc ĐƯỢC XEN KẼ (có thể chỉ trên 1 lõi CPU).
//	PARALLEL    = nhiều việc chạy THẬT SỰ ĐỒNG THỜI (cần nhiều hơn 1 lõi).
//
// Goroutine cho bạn concurrency; parallelism phụ thuộc GOMAXPROCS và số lõi thật.
func ConcurrencyVsParallelism() {
	runtime := "GOMAXPROCS hiện tại = số lõi CPU Go được phép dùng"
	_ = runtime

	jobs := make(chan int, 4)
	done := make(chan struct{})

	// 3 worker — "concurrent" bất kể máy có bao nhiêu lõi.
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for j := range jobs {
				fmt.Printf("  worker %d xử lý job %d\n", id, j)
				time.Sleep(50 * time.Millisecond) // giả lập công việc
			}
			done <- struct{}{}
		}(w)
	}

	for j := 1; j <= 4; j++ {
		jobs <- j
	}
	close(jobs)
	for w := 1; w <= 3; w++ {
		<-done
	}
}

// GoroutineLeak: goroutine này block VĨNH VIỄN ở dòng `ch <- 1`
// vì không còn ai đọc channel nữa → bộ nhớ + slot scheduler bị giữ mãi.
// Đây là lỗi rò rỉ goroutine kinh điển: không có cancellation / không có deadline.
func GoroutineLeak() {
	ch := make(chan int)

	go func() {
		ch <- 1 // block mãi: main không còn đọc ch → không bao giờ thoát
		fmt.Println("dòng này không bao giờ được in")
	}()

	_ = <-ch // main chỉ đọc MỘT lần rồi bỏ đi
	fmt.Println("main thoát, goroutine kia vẫn còn bị kẹt")
}

// DataRace: counter++ là read-modify-write, 2 goroutine cùng làm thì MẤT CẬT NHẬT.
// Chạy với `go run -race .` để toolchain phát hiện ra race condition.
func DataRace() {
	const n = 1000
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // KHÔNG an toàn: không có lock/atomic
		}()
	}
	wg.Wait()
	fmt.Printf("counter = %d (mong đợi %d) — nếu nhỏ hơn là do data race\n", counter, n)
}

// LoopCapture: mỗi iteration phải có biến của RIÊNG nó.
// Từ Go 1.22+, `for i := range n` đã tạo biến i MỚI mỗi iteration
// nên bug kinh điển "n goroutine đều thấy i = n" không còn tái hiện.
// VẪN CẦN chú ý khi capture biến KHÁC (vd: `v := v` cho biến ngoài loop).
func LoopCapture() {
	var wg sync.WaitGroup // dùng ở đây chỉ để ĐỢI in cho đủ, không dạy sync
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  goroutine thấy i =", i)
		}()
	}
	wg.Wait()
}
