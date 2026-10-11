package main

import "fmt"

func Slice() {
	arr := []int{1, 3, 5, 5, 767}
	s := []int{1, 2, 3}
	//make([]Type, length, capacity)
	s2 := make([]int, 3, 5)
	s3 := arr[1:4]
	// Thay đổi s3 thì sẽ thay đổi cả phần tử đó tại underlaying data
	s3[0] = 123
	s3 = append(s3, 4)
	fmt.Println("Test arr,s,s2", arr, s, s2)
	fmt.Println("Test s3", s3)

	// Để tránh thay đổi ở underlaying data có thể dụng copy, chỉ dùng với slice
	fmt.Println("before arr", arr)
	s4 := make([]int, 2, 4)
	copy(s4, arr)
	fmt.Println("after arr", arr)

	s5 := arr[1:4:4]
	fmt.Println("Test s5", arr, s5)

	// nill slice & emtpy slice
	var nillSlice [4]int   // nill
	emptySlice := [4]int{} // empty
	fmt.Println("Check nill và emtpy", nillSlice, emptySlice)
}
