package main

import "fmt"

func Struct() {
	type Address struct {
		City    string
		Country string
	}
	type User struct {
		ID      int
		Name    string
		Age     int
		Address Address
	}
	trungLe := User{
		ID:   1,
		Name: "Trung",
		Age:  24,
		Address: Address{
			City:    "Thanh hoa",
			Country: "Viet Nam",
		},
	}
	// Struct embedding
	type Admin struct {
		User
		Permission []string
	}
	fmt.Println(trungLe.Address.City)
	trungLe2 := User{
		ID:   1,
		Name: "Trung",
		Age:  24,
		Address: Address{
			City:    "Thanh hoa",
			Country: "Viet Nam",
		},
	}
	fmt.Println(trungLe == trungLe2)
}
