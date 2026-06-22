package main

import (
	"fmt"
)

func PrintNumber(nums [5]int) {
	fmt.Println("nums", nums)
}

func Array() {
	// 1. Nhóm array ( Array) : Mảng có kích thước cố định và lưu các phần tử cùng cùng kiểu dữ liệu
	var a [3]int               // [0,0,0]
	b := [7]int{1, 2, 3, 4, 5} // [1,2,3,4,5,0,0]
	c := [...]string{"a", "b", "c"}
	students := [...]string{"Lê Văn Trung", "Hoàng Việt Hùng"}
	for _, student := range students {
		fmt.Println("Student in ", student, len(students))
	}

	fmt.Println("Array values ", a, b, c[2])
	fmt.Println("Array length ", len(a), len(b), len(c))
	a1 := [5]int{1, 2, 3, 4, 5}
	// b2 := [3]int{1, 2, 3}
	PrintNumber(a1) // OK
	// Khi dùng Array, lenght là 1 phần của type nên báo error, slice thì bình thường
	// PrintNumber(b2) // Error
}
