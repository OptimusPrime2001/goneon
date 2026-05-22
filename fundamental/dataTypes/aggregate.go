package main

import "fmt"

// Aggregate types là kiểu tổng hợp gom nhóm
func AggregateTypes() {
	// 1. Nhóm array ( Array) : Mảng có kích thước cố định và lưu các phần tử cùng cùng kiểu dữ liệu
	// Array()

	// 2. Nhóm slice ( Slice) : Mảng không cố định
	Slice()

	// 3. Nhóm map ( Map) : Lưu các phần tử khác nhau theo cặp key-value
	// Map()

	// 4. Nhóm struct ( Struct)
	// Struct()
}
func Array() {
	// 1. Nhóm array ( Array) : Mảng có kích thước cố định và lưu các phần tử cùng cùng kiểu dữ liệu
	var a [3]int               // [0,0,0]
	b := [7]int{1, 2, 3, 4, 5} // [1,2,3,4,5,0,0]
	c := [...]string{"a", "b", "c"}
	fmt.Println("Array values ", a, b, c[2])
	fmt.Println("Array length ", len(a), len(b), len(c))
}
func Slice() {
	// 2. Nhóm slice ( Slice) : Mảng không cố định
	// len là số phần tử Slice đang dùng
	// cap là max số phần tử tối đa slice có thể dùng trước khi Go phải cấp phát memory mới.
	nums := []int{1, 2, 3}
	fmt.Println(len(nums), cap(nums))
	nums = append(nums, 4, 5)

	fmt.Println(len(nums), cap(nums))
	b := make([]string, 0, 5) // len 0, cap 5

	// Slice from array
	arr := [5]int{1, 2, 3, 4, 5}
	s := arr[1:3] // len 2, cap 4. Cap bắt đầu từ index 1 đến cuối array
	s = append(s, 12)
	fmt.Println("Slice from array ", b, s, arr)

}
func Map() {

}
func Struct() {

}

/*
>>> Các lưu ý về array trong go
1: Mảng khai báo phải có kích thước cố định, ko thể thay đổi kích thước sau khi khai báo
2: Nếu muốn thao tác thêm, sửa xoá thì cần dùng slice ( Slice )
3: Slice linh hoạt hơn array nên được sử dụng nhiều hơn
*/
