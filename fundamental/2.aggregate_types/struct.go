package main

import "fmt"

type Addres struct {
	Street string
	City   string
}
type Person struct {
	Name    string
	Age     int
	Phone   string
	Address Addres
}

var defaultPerson = Person{
	Name:  "Trung",
	Age:   20,
	Phone: "0987654321",
	Address: Addres{
		Street: "123 Main St",
		City:   "New York",
	},
}

func Basic_Struct() {
	trung := defaultPerson
	nam := &trung
	nam.Name = "Nam"
	nam.Address.Street = "456 Elm St"
	fmt.Printf("Trung: %v, \nNam: %v\n", trung, *nam)
}
func (p Person) GetInfo() {
	fmt.Println("Get info : ", p)
}
func (p *Person) SetInfo() {
	p.Name = "Nguyen Thao An"
	p.Age = 24
	p.Phone = "0987654321"
	p.Address.Street = "Đông Anh Hà Nội"
	p.Address.City = "Hà Nội"
	fmt.Printf("Thao An info : %p\n", p)
}
func Struct() {
	Basic_Struct()
	newPerson := Person{}
	// Value Receiver
	newPerson.GetInfo()

	// Pointer Receiver
	newPerson.SetInfo()
}
