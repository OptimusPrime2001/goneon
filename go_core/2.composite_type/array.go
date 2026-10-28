package main

import "fmt"

func Array() {
	var numbers [5]int
	var numbers2 = []int{12, 2, 33, 74, 25}
	var test = numbers2
	// Array là value type, ko phải reference type nên test thay đổi thế nào cũng ko ảnh hưởng đế numbers2
	test[0] = 1
	fmt.Println("Test ", test)
	numbers2[0] = 10
	fmt.Println("numbers2 ", numbers2)
	matrix := [][]int{
		{1, 2, 3, 4},
		{2, 3, 54, 6},
	}

	fmt.Println("Numbers array", numbers, matrix)
}
