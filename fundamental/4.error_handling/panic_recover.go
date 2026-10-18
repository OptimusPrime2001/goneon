package main

import "fmt"

// PanicDemo minh họa: panic dừng hàm, nhưng các defer ĐÃ đăng ký vẫn chạy (LIFO).
func PanicDemo() {
	fmt.Println("Start")
	defer fmt.Println("Cleanup 1")
	defer fmt.Println("Cleanup 2") // chạy TRƯỚC Cleanup 1 (LIFO)

	panic("something went wrong")
	// Dòng dưới KHÔNG bao giờ được in ra.
	fmt.Println("End")
}

// RecoverDemo minh họa: recover() trong defer bắt panic, chương trình sống tiếp.
func RecoverDemo() {
	defer func() {
		// recover() trả về giá trị panic nếu có, ngược lại nil.
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()

	fmt.Println("Before panic")
	panic("something went wrong")

	// KHÔNG được thực thi.
	fmt.Println("After panic")
}

// DivideSafe bọc hàm prone-to-panic rồi recover, trả về error bình thường.
// Đây là MÔ HÌNH chuẩn: panic → recover → chuyển thành error.
func DivideSafe(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DivideSafe: %v", r) // chuyển panic thành error
		}
	}()
	return a / b, nil // panic runtime nếu b == 0
}
