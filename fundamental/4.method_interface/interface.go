package main

import "fmt"

// --- 1. Interface basic ---
// KHÔNG có keyword "implements" — implement IMPLICIT: chỉ cần có đủ method.

type Shape interface {
	Area() float64
}

type Circle struct{ Radius float64 }

// Circle implement Shape vì có Area() float64 — không cần khai báo gì thêm.
func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

type Triangle struct{ Width, Height float64 }

func (t Triangle) Area() float64 {
	return 0.5 * t.Width * t.Height
}

// PrintArea nhận MỌI type nào implement Shape.
func PrintArea(s Shape) {
	fmt.Println("Area:", s.Area())
}

// --- 2. Empty interface ---
// any = interface{} — chấp nhận MỌI giá trị.

func printAnything(v any) {
	fmt.Println("any:", v)
}

// --- 3. Type assertion ---
func typeAssertion(v any) {
	// Comma-ok: an toàn, không panic
	s, ok := v.(string)
	if ok {
		fmt.Println("assert string:", s)
	} else {
		fmt.Println("assert: không phải string")
	}
}

// --- 4. Type switch ---
func typeName(v any) {
	switch v.(type) {
	case int:
		fmt.Println("type: int")
	case string:
		fmt.Println("type: string")
	case Circle:
		fmt.Println("type: Circle")
	default:
		fmt.Println("type: unknown")
	}
}

// --- 5. Embedding interface ---
// Gộp 2 interface nhỏ thành lớn (mô hình io.ReadWriter).

type Reader interface{ Read() string }
type Writer interface{ Write(s string) }

// ReadWriter = Reader + Writer — embedding, không viết lại method.
type ReadWriter interface {
	Reader
	Writer
}

type File struct{ name string }

func (f File) Read() string   { return "data from " + f.name }
func (f File) Write(s string) { fmt.Println("write", s, "to", f.name) }

func Interfaces() {
	// Implicit implementation
	PrintArea(Circle{Radius: 2})
	PrintArea(Triangle{Width: 4, Height: 3})

	// Empty interface
	printAnything(42)
	printAnything("hello")

	// Type assertion
	typeAssertion("xin chao")
	typeAssertion(3.14)

	// Type switch
	typeName(1)
	typeName("hi")
	typeName(Circle{})

	// Embedding interface
	var rw ReadWriter = File{name: "note.txt"}
	fmt.Println(rw.Read())
	rw.Write("hello")
}
