package main

import (
	"fmt"
	"slices"
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
	InitSlice()
	CopySlice()
	NillAndEmpty()
	SliceOperations(originalArray)
	DeleteItem()
	InsertSlice()
}
func InitSlice() {
	s1 := []string{"BTS", "Big Bang", "Cortis"}
	var s2 = []int{1, 3, 4, 5}
	s3 := make([]int, 3, 5)
	fmt.Println("======== InitSlice=========")
	fmt.Printf("Init Slice s1 %v, len: %d, cap: %d \n", s1, len(s1), cap(s1))
	fmt.Printf("Init Slice s2 %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))
	fmt.Printf("Init Slice s3 %v, len: %d, cap: %d \n", s3, len(s3), cap(s3))
}
func CopySlice() {
	// Tạo slice mới từ array (Ko ảnh hưởng slice gốc)
	// s := slices.Clone(arr[1:3])

	// s := make([]int, len(arr[1:3]))
	// copy(s, arr[1:3])

	// s := append([]int{}, arr[1:3]...)
	src := []int{1, 2, 3, 4, 5}
	dst1 := make([]int, len(src))

	copy(dst1, src)
	dst1[3] = 10
	dst2 := slices.Clone(src)
	dst2[3] = 29
	dst3 := append([]int{}, src...)
	dst3[3] = 111
	fmt.Println("======== CopySlice =========")
	fmt.Printf("Slice original %v \n", src)
	fmt.Printf("Slice original %v \n", dst1)
	fmt.Printf("Slice original %v \n", dst2)
	fmt.Printf("Slice original %v \n", dst3)
}
func NillAndEmpty() {
	fmt.Println("======== NillAndEmpty =========")
	var nillSlice []int
	fmt.Printf("Check nill slice %v, len: %d, cap: %d \n", nillSlice, len(nillSlice), cap(nillSlice))

	emptySlice := []int{}
	fmt.Printf("Check nill slice %v, len: %d, cap: %d", emptySlice, len(emptySlice), cap(emptySlice))
}
func SliceOperations(slice []int) {
	fmt.Println("\n======== SliceOperations =========")
	original := slice
	fmt.Println("Slice gốc :", original)
	fmt.Println("original[1:4] : ", original[1:4]) // [2 3 4]
	fmt.Println("original[:3] : ", original[:3])   // [1 2 3]
	fmt.Println("original[2:] : ", original[2:])   // [3 4 5]
}
func DeleteItem() {
	fmt.Println("\n======== DeleteItem =========")
	s := []int{0, 1, 2, 3, 4}

	// xóa index 1
	s = append(s[:2], s[3:]...)
	fmt.Println("Slice deleted: ", s)

}
func InsertSlice() {
	fmt.Println("\n======== InsertSlice=========")
	slice := []int{1, 3, 5}
	fmt.Println("After insert slice", slice)
	slice = append(slice[:1], append([]int{2}, slice[1:]...)...)
	fmt.Println("Before insert Slice", slice)
}
