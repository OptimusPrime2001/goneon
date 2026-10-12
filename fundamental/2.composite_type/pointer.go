package main

import "fmt"

func basic() {
	i, j := 42, 1701
	pointer := &i
	fmt.Println("Check p", pointer)
	*pointer = 21
	fmt.Println("Check i", i)

	pointer = &j
	*pointer = *pointer / 37
	fmt.Println("Check j", j)
}
func withStruct() {
	type AppUser struct {
		fname string
		lname string
		age   int
	}
	appUser := AppUser{
		fname: "Hasin",
		lname: "Shekih",
		age:   35,
	}
	pointer := &appUser
	pointer.age = 15
	fmt.Println("Check app user", appUser)
}
func mapsAndSlice() {
	person := map[string]int{
		"age": 25,
	}
	change := func(m map[string]int) {
		// Sửa giá trị
		m["age"] = 200
		// Gán lại biến
		m = map[string]int{
			"age": 31,
		}
	}
	change(person)
	fmt.Println("Check pointer with map", person["age"])

	numbers := []int{1, 2, 3, 4, 5}
	changeSlice := func(slice []int) {
		slice[1] = 200
	}
	changeSlice(numbers)

	fmt.Println("Check pointer with slice", numbers)
}

// Trường hợp đặc biệt khi dùng Slice với Append

func add(s []int) {
	s = append(s, 4)
	// return append(s, 4)
}
func sliceWithAppend() {
	numbers := []int{1, 2, 3}

	add(numbers)

	fmt.Println("Slice with append", numbers) // [1 2 3 4]
}
func Pointer() {
	basic()
	withStruct()
	mapsAndSlice()
	sliceWithAppend()

}
