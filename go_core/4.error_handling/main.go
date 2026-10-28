package main

import "fmt"

func main() {
	fmt.Println("=== Error Basic ===")
	ErrorBasic()

	fmt.Println("\n=== Wrapping ===")
	Wrapping()

	fmt.Println("\n=== Sentinel ===")
	Sentinel()

	fmt.Println("\n=== Recover ===")
	RecoverDemo()

	if r, err := DivideSafe(10, 0); err != nil {
		fmt.Println("DivideSafe lỗi:", err)
	} else {
		fmt.Println("DivideSafe =", r)
	}

	fmt.Println("\n=== Stack Trace ===")
	StackTrace()

	// PanicDemo cố ý panic KHÔNG recover → crash chương trình.
	// Đặt CUỐI CÙNG để các phần trên đã in xong output.
	fmt.Println("\n=== Panic (crash chương trình) ===")
	PanicDemo()
}
