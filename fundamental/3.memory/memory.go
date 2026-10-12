package main

import "fmt"

// escapeToHeap: biến local x bị return dưới dạng pointer → compiler đẩy x lên HEAP.
// (Kiểm chứng: go build -gcflags="-m" → "moved to heap: x")
func escapeToHeap() *int {
	x := 42
	return &x
}

// stayInStack: biến local x chỉ dùng trong hàm, không escape → nằm STACK.
func stayInStack() {
	x := 42
	_ = x
}

// closureCapture: closure "nhớ" biến n → n chuyển lên HEAP để sống lâu hơn hàm ngoài.
func closureCapture() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// boxing: biến int bị "đóng gói" vào interface{} → escape ra HEAP.
func boxing() {
	var v any = 42
	fmt.Println(v)
}

func Memory() {
	stayInStack()

	p := escapeToHeap()
	fmt.Println("escapeToHeap:", *p)

	counter := closureCapture()
	fmt.Println("closureCapture:", counter(), counter())

	boxing()
}
