package main

import (
	"fmt"
	"runtime/debug"
)

// riskyFunc cố ý panic để demo stack trace.
func riskyFunc() {
	panic("đã xảy ra lỗi nghiêm trọng")
}

// middleLayer chỉ là frame trung gian để trace dài hơn.
func middleLayer() {
	riskyFunc()
}

func StackTrace() {
	// recover() phải nằm TRONG hàm defer mới có tác dụng.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
			fmt.Println("---- Stack trace (debug.Stack) ----")
			// debug.Stack() trả về stack hiện tại của GOROUTINE tại thời điểm gọi.
			// Dòng đầu tiên là "goroutine N [running]:", các dòng sau là các frame.
			fmt.Println(string(debug.Stack()))
		}
	}()

	middleLayer()
}
