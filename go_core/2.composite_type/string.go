package main

import "fmt"

func String() {
	const sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
	fmt.Printf("% x\n", sample)
	fmt.Printf("%q\n", sample)

	for i := 0; i < len(sample); i++ {
		fmt.Printf("%x ", sample[i])
	}
	var b byte = 0x41 // or 65

	fmt.Printf("\n%d\n", b) // 65  — decimal
	fmt.Printf("%x\n", b)   // 41  — hexadecimal
	fmt.Printf("%b\n", b)   // 1000001 — binary
	fmt.Printf("%c\n", b)   // A   — character
}
