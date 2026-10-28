package main

import "fmt"

func Maps() {
	// map[KeyType]ValueType
	// Với nill map thì chỉ có thể read, nếu ghi sẽ gây lỗi panic
	var m map[string]int
	m = make(map[string]int)
	m["route"] = 66
	i, ok := m["route"]
	fmt.Println(i, ok)
	n := len(m)
	fmt.Println(n)
	delete(m, "route")
	clear(m)
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}

	giaTraiCay := map[string]int{
		"Apple":  20,
		"Orange": 25,
		"Banana": 31,
	}
	for k := range giaTraiCay {
		fmt.Println("Key:", k, "Value:", giaTraiCay[k])
	}
}
