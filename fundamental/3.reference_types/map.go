package main

import "fmt"

func Map() {
	giaTraiCay := map[string]int{
		"Apple":  20,
		"Orange": 25,
		"Banana": 31,
	}
	giaTraiCay["Apple"] = 22
	giaTraiCay["Mango"] = 26
	gia, exists := giaTraiCay["Apple"]
	if exists {
		fmt.Println("Apple price is ", gia)
	}
	delete(giaTraiCay, "Cam") // Xóa Cam khỏi map
	for k, v := range giaTraiCay {
		fmt.Printf("Qua %s co gia la %d\n", k, v)
	}

}
