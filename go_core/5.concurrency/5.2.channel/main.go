package main

import (
	"fmt"
	"sync"
	"time"
)

func UnbufferedChannel() {
	ch := make(chan int)
	go func() {
		fmt.Println("Trước khi gửi dữ liệu")
		// Chỗ này sẽ bị block lại, ko chạy xuống dưới nếu chưa có receiver <- ch
		ch <- 100
		fmt.Println("Đã gửi dữ liệu")
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("Nhận được dữ liệu: ", <-ch)
}
func BufferedChannel() {
	ch2 := make(chan int, 3)
	ch2 <- 11
	ch2 <- 12
	ch2 <- 20
	close(ch2)
	fmt.Println("Nhận lần đầu", <-ch2)
	// Nhận và in ra các dữ liệu bên trong buffer
	// for value := range ch2 {
	// 	fmt.Println("Các giá trị bên trong buffer", value)
	// }
	// Gây panic nếu gửi dữ liệu vào channel đã close
	value, ok := <-ch2
	if !ok {
		fmt.Println("Channel đã hết dữ liệu", value)
		ch2 <- 30
	} else {
		fmt.Println("Nhận lần hai", value)
		// len(channel) Số tra giá trị đang chờ dc nhận
		fmt.Println("Số giá trị đang chờ nhận trong buffer", len(ch2))

	}
	fmt.Println("Nhận lần ba", <-ch2)
}
func SelectStatement() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(time.Second)
		ch1 <- "From channel 1"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "From channel 2"
	}()

	select {
	case msg := <-ch1:
		fmt.Println("Nhận từ channel 1", msg)
	case msg := <-ch2:
		fmt.Println("Nhận từ channel 2", msg)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
		// default:
		// 	fmt.Println("No data available")
	}
}

func BAI1() {
	ch := make(chan int)
	go func() {
		ch <- 100
	}()

	value := <-ch
	fmt.Println("Nhận value", value)
}
func BAI2() {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	for v := range ch {
		fmt.Println(v)
	}
}
func BAI3() {
	ch := make(chan int, 2)
	go func() {
		for i := 1; i < 4; i++ {
			ch <- i
			fmt.Println("Sent: ", i)
		}
	}()
	v1 := <-ch
	fmt.Println("Received: ", v1)
	v2 := <-ch
	fmt.Println("Received: ", v2)
	v3 := <-ch
	fmt.Println("Received: ", v3)
}

func BAI6() {
	ch := make(chan int, 3)
	go func() {
		for i := 1; i < 6; i++ {
			ch <- i
			fmt.Println("Produced: ", i)
		}
		close(ch)
	}()
	func() {
		for v := range ch {
			fmt.Println("Consumed :", v)
		}
	}()
}
func BAI7() {
	ch := make(chan int, 3)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i < 6; i++ {
			ch <- i
			fmt.Println("Produced: ", i)
		}

	}()
	go func() {
		// Đảm bảo gọi Done khi hàm kết dùng, dùng defer với mục đích vậy
		defer wg.Done()
		for i := 101; i < 106; i++ {
			ch <- i
			fmt.Println("Produced: ", i)
		}

	}()
	go func() {
		// Chặn, block goroutine gọi nó cho đến khi bộ đếm = 0
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println("Consumed:", v)
	}

}
func BAI8() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		ch1 <- 10
		fmt.Println("Sent ch1: ", 10)
	}()
	go func() {
		ch2 <- 20
		fmt.Println("Sent ch2: ", 20)
	}()

	for i := 0; i < 2; i++ {
		select {
		case v1 := <-ch1:
			fmt.Println("Received from ch1:", v1)
		case v2 := <-ch2:
			fmt.Println("Received from ch2:", v2)
		}
	}
}
func BAI9() {
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		fmt.Println("Received: ", v)
	default:
		fmt.Println("No data available")
	}
	ch <- 99
	select {
	case v := <-ch:
		fmt.Println("Received: ", v)
	default:
		fmt.Println("Nove data available")
	}
}
func BAI10() {
	jobs := make(chan int, 2)
	results := make(chan int, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i < 4; i++ {
			jobs <- i
			fmt.Println("Job sent: ", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 101; i < 104; i++ {
			jobs <- i
			fmt.Println("Job sent: ", i)
		}
	}()
	go func() {
		wg.Wait()
		close(jobs)
	}()

	var wg2 sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for job := range jobs {
				result := job * 2
				results <- result
				fmt.Println("Result sent: ", result)
			}
		}()
	}
	go func() {
		wg2.Wait()
		close(results)
	}()
	received := 0
	for received < 6 {
		v, ok := <-results
		if !ok {
			fmt.Println("Results channel closed")
			return
		}
		fmt.Println("Final result:", v)
		received++
	}

}
func main() {
	// UnbufferedChannel()
	// BufferedChannel()
	// SelectStatement()
	// BAI1()
	// BAI2()
	// BAI6()
	// BAI7()
	BAI10()
}
