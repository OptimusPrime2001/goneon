package main

import "fmt"

func BasicTypes() {
	// 1. Nhóm số nguyên ( Integer)
	// 1.1: Số nguyên có dấu gồm cả số âm và số dương (Signed Integer) : int, int8, int16, int32, int64
	var count int8 = 123
	var score int = -23432 // với int thì kích thước phụ thuộc vào OS là 32 bit hay 64 bit
	fmt.Println("count int8", count)
	fmt.Println("score int", score)
	// 1.2: Số nguyên không có dấu chỉ gồm cả số dương (Unsigned Integer) : uint, uint8, uint16, uint32, uint64
	var age uint8 = 25
	fmt.Println("age unit8", age)
	// 1.3: Số nguyên đặc biệt : byte, rune
	var b byte = 100
	var r rune = 'a'
	fmt.Println("b byte", b)
	fmt.Println("r rune", r)

	// 2. Nhóm số thực ( Float)
	// 2.1: Số thực có dấu (Signed Float) : float32, float64
	var height float32 = 1.85
	fmt.Println("height float32", height)
	var weight float64 = 73.2134
	fmt.Println("weight float64", weight)

	// 3. Nhóm bool ( Boolean)
	var isActive bool = true
	fmt.Println("isActive bool", isActive)

	// 4. Nhóm string ( String)
	name := "Le Van Trung"
	msg := "Xin chao.\n tôi là " + name
	fmt.Println("name string", name)
	fmt.Println("msg string", msg)
	//
}
