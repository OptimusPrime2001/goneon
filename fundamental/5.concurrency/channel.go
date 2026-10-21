package main

import (
	"fmt"
	"time"
)

// UnbufferedDemo: hand-off đồng bộ — người gửi CHỜ có người nhận.
func UnbufferedDemo() {
	ch := make(chan int) // cap = 0 → đồng bộ

	go func() {
		fmt.Println("  sender: chuẩn bị gửi")
		ch <- 42 // đứng chờ tới khi main đọc
		fmt.Println("  sender: đã gửi xong, chạy tiếp được rồi")
	}()

	v := <-ch // main chờ tới khi có giá trị
	fmt.Println("  receiver: nhận được", v)
}

// BufferedDemo: channel có hàng đợi — gửi KHÔNG cần người nhận ngay,
// cho tới khi đầy (len == cap) thì người gửi block.
func BufferedDemo() {
	ch := make(chan int, 2) // cap = 2

	ch <- 1
	ch <- 2
	// ch <- 3  // nếu uncomment → block: hàng đợi đã đầy, không ai đọc

	fmt.Println("  len =", len(ch), "cap =", cap(ch))

	fmt.Println("  nhận:", <-ch)
	fmt.Println("  sau khi nhận 1 phần tử: len =", len(ch), "cap =", cap(ch))
}

// DirectionalDemo: tham số channel chỉ định HƯỚNG dữ liệu được phép đi.
// → ràng buộc contract: hàm producer chỉ được GỬI, consumer chỉ được NHẬN.
func producer(out chan<- int) { // chỉ được gửi
	for i := 1; i <= 3; i++ {
		out <- i * 10
	}
	close(out) // người GỬI là người close
}

func consumer(in <-chan int) { // chỉ được nhận
	for v := range in { // range tự dừng khi channel bị close
		fmt.Println("  consumer nhận:", v)
	}
}

func DirectionalDemo() {
	ch := make(chan int)
	go producer(ch)
	consumer(ch) // main đóng vai consumer, block tới khi producer close
}

// CloseRangeDemo: close + range — "đóng cửa, hết hàng thì hết vòng lặp".
func CloseRangeDemo() {
	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"
	close(ch)

	for v := range ch { // lặp tới khi channel đóng rồi mới thoát
		fmt.Println("  range nhận:", v)
	}
	fmt.Println("  range kết thúc: channel đã đóng")
}

// CommaOkDemo: phân biệt "channel đóng" với "giá trị 0 hợp lệ".
func CommaOkDemo() {
	ch := make(chan int, 1)
	ch <- 7
	close(ch)

	v, ok := <-ch
	fmt.Printf("  v=%d ok=%v  (ok=true → còn dữ liệu)\n", v, ok)

	v, ok = <-ch
	fmt.Printf("  v=%d ok=%v  (ok=false → channel đóng, v là zero value)\n", v, ok)
}

// SelectDemo: đọc từ NHIỀU channel — case nào SẴN SÀNG trước thì chọn.
// Nếu nhiều case sẵn sàng → Go chọn NGẪU NHIÊN (không theo thứ tự).
func SelectDemo() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() { time.Sleep(50 * time.Millisecond); c1 <- "từ c1" }()
	go func() { time.Sleep(100 * time.Millisecond); c2 <- "từ c2" }()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-c1:
			fmt.Println("  nhận:", msg)
		case msg := <-c2:
			fmt.Println("  nhận:", msg)
		}
	}
}

// SelectDefaultDemo: default → KHÔNG block, chạy ngay khi chưa case nào sẵn sàng.
// Đây là cách làm channel "non-blocking".
func SelectDefaultDemo() {
	ch := make(chan int, 1) // buffered, nhưng để TRỐNG

	select {
	case v := <-ch:
		fmt.Println("  nhận:", v)
	default:
		fmt.Println("  không có gì sẵn → chạy default NGAY, không chờ")
	}

	ch <- 99
	select {
	case v := <-ch:
		fmt.Println("  có dữ liệu → nhận:", v)
	default:
		fmt.Println("  default (sẽ không chạy)")
	}
}

// SelectTimeoutDemo: time.After trả channel "báo sau X" — nếu không nhận
// kịp thì chọn nhánh timeout.
func SelectTimeoutDemo() {
	ch := make(chan string) // không bao giờ có dữ liệu

	select {
	case v := <-ch:
		fmt.Println("  nhận:", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  timeout: không có dữ liệu sau 100ms")
	}
}

// NilSelectDemo: case với channel nil KHÔNG BAO GIỜ sẵn sàng → bị bỏ qua.
// Dùng để "tắt" 1 nhánh select có điều kiện.
func NilSelectDemo() {
	var ch chan int // nil — không có case nào chạy được từ nó

	select {
	case v := <-ch:
		fmt.Println("  không bao giờ chạy:", v) // KHÔNG chạy
	case <-time.After(50 * time.Millisecond):
		fmt.Println("  case nil bị bỏ qua → timeout thắng")
	}
}

func Channel() {
	UnbufferedDemo()
	BufferedDemo()
}
