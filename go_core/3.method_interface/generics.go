package main

import (
	"cmp"
	"fmt"
)

// --- 1. Why generics ---
// Trước Go 1.18, mỗi kiểu cần1 hàm riêng:
func SumInt(nums ...int) int {
	var total int
	for _, n := range nums {
		total += n
	}
	return total
}

func SumFloat(nums ...float64) float64 {
	var total float64
	for _, n := range nums {
		total += n
	}
	return total
}

// → Generics gộp tất cả thành 1 hàm Sum[T Number] bên dưới.

// --- 2. Custom constraint: TYPE UNION + tilde (~) ---
// ~int = int HOẶC kiểu nào có underlying type là int (vd: type Celsius int)

type Number interface {
	~int | ~int64 | ~float64
}

// --- 3. Generic function ---
// T là TYPE PARAMETER — Sum(1,2,3) và Sum(1.5,2.5) dùng chung1 hàm.

func Sum[T Number](nums ...T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

// Max dùng constraint có sẵn của stdlib: cmp.Ordered (Go 1.21+)
// (ordered = <, >, <=, >= hoạt động: số, string)

func Max[T cmp.Ordered](first T, rest ...T) T {
	m := first
	for _, v := range rest {
		if v > m {
			m = v
		}
	}
	return m
}

// --- 4. Generic type (struct) ---
// Stack[T] —1 kiểu cho mọi kiểu phần tử.

type Stack[T any] struct {
	items []T
}

// Method của generic type cũng khai báo [T any].
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // zero value của T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func Generics() {
	// --- Type inference: compiler tự suy luận T từ tham số ---
	fmt.Println("Sum int:", Sum(1, 2, 3, 4))
	fmt.Println("Sum float:", Sum(1.5, 2.5))

	// Type argument tường minh
	fmt.Println("Sum explicit:", Sum[int](1, 2, 3))

	fmt.Println("Max:", Max(5, 2, 9, 4))
	fmt.Println("Max string:", Max("banana", "apple", "cherry"))

	// Generic stack
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	v, ok := s.Pop()
	fmt.Println("Pop:", v, ok)
}
