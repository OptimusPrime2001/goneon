package main

import (
	"fmt"
)

/*
>>> Các lưu ý về Slice trong go
1: Slice chỉ là “view” trỏ tới underlying array, sửa slice con thì slice cha cũng bị ảnh hưởng
2: append() dùng chung array nếu cap còn đủ, tạo array mới nếu cap cũ không đủ
3: len() là số phần tử Slice đang dùng
4: cap() là max số phần tử tối đa slice có thể dùng trước khi Go phải cấp phát memory mới.
5: Những cách tạo slice mới mà ko ảnh hưởng slice gốc
  - append
  - make + copy
  - slices.Clone()
*/
func Slice() {
	originalArray := []int{1, 2, 3, 4, 5, 6}
	newSlice := originalArray[1:4]
	fmt.Println("Init Slice ", len(originalArray), cap(originalArray))
	nums := append(originalArray, 4, 5)

	fmt.Println("Append Slice ", len(nums), cap(nums))
	b := make([]int, 2, 5) // len 2, cap 5

	// Slice from array
	arr := [5]int{1, 2, 3, 4, 5}

	// Tạo slice mới từ array (Ko ảnh hưởng slice gốc)
	// s := slices.Clone(arr[1:3])

	// s := make([]int, len(arr[1:3]))
	// copy(s, arr[1:3])

	// s := append([]int{}, arr[1:3]...)

	// Tạo slice mới từ array (Ảnh hưởng slice gốc)
	s := arr[1:3] // len 2, cap 4. Cap bắt đầu từ index 1 đến cuối array
	// Cả s và arr gốc đều thay đổi khi thay đổi s
	s[0] = 9
	s = append(s, 12)
	empty := []int{}
	fmt.Println("Empty slice", empty)
	fmt.Println("Slice from array ", b, s, arr)
}
func CopySlice() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))

	copy(dst, src)
}
func NillAndEmpty() {
	var nillSlice []int
	fmt.Printf("Check nill slice %v, len: %d, cap: %d", nillSlice, len(nillSlice), cap(nillSlice))

	emptySlice := []int{}
	fmt.Printf("Check nill slice %v, len: %d, cap: %d", emptySlice, len(emptySlice), cap(emptySlice))
}
func SliceOperations(slice []int) {
	original := slice
	fmt.Println("Slice gốc :", original)
	fmt.Println("original[1:4] : ", original[1:4]) // [2 3 4]
	fmt.Println("original[:3] : ", original[:3])   // [1 2 3]
	fmt.Println("original[2:] : ", original[2:])   // [3 4 5]
}
func DeleteItem() {

}
func InsertSlice() {
	slice := []int{1, 3, 5}
	fmt.Println("After insert slice", slice)
	slice = append(slice[:1], append([]int{2}, slice[1:]...)...)
	fmt.Println("Before insert Slice", slice)
}
