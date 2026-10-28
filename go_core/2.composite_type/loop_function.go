package main

import "fmt"

func Loops() {
	sum := 0
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue
		}
		if i == 8 {
			goto AnThao
		}
		fmt.Printf("%d ", i)
		sum += i
	}
	fmt.Print(sum)
AnThao:
	fmt.Println("Chúc mừng Thảo An")

	switch sum {
	case 20:
		goto AnThao
	}
}

func test(x int, y int) int {
	return x + y
}
func an() (bool, int) {
	return true, 2
}
func variadic(numbers ...int) {
	fmt.Println(numbers)
}
func clouse() func() int {
	x := 10

	return func() int {
		x++
		return x
	}
}

// Named return values
func calculate(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b

	return sum, diff // hoặc return trống cũng có thể hiểu
}
func example() (n int, ok bool, message string) {
	return
}

// Call by value: Go LUÔN copy tham số khi gọi hàm.
// int / struct  -> bản copy, sửa trong hàm không ảnh hưởng biến gốc.
// map / slice / pointer -> copy "con trỏ", dùng chung dữ liệu gốc.
type NhanVien struct {
	Luong int
}

func suaNhanVien(luong int, nv NhanVien, giaTri map[string]int) {
	luong = 999       // chỉ sửa bản copy -> biến gốc không đổi
	nv.Luong = 999    // struct cũng là kiểu giá trị -> không đổi
	giaTri["x"] = 999 // map là reference -> biến gốc ĐỔI
}

func CallByValue() {
	luong := 10
	nv := NhanVien{Luong: 10}
	giaTri := map[string]int{"x": 1}

	suaNhanVien(luong, nv, giaTri)

	fmt.Println("luong:", luong)       // 10   không đổi
	fmt.Println("nv.Luong:", nv.Luong) // 10   không đổi
	fmt.Println("giaTri:", giaTri)     // map[x:999]  ĐỔI
}

func Functions() {
	test(2, 4)
	// Variadic Functions
	variadic(1, 34, 54, 5, 65)
	an()

	// anonymous functon
	add := func(a, b int) int {
		return a + b
	}
	fmt.Println(add(3, 4))

	counter := clouse()

	fmt.Println(counter()) // 11
	fmt.Println(counter()) // 12
	fmt.Println(counter()) // 13

	sum, diff := calculate(10, 3)
	fmt.Println(sum, diff) // 13 7

	CallByValue()
}
