package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Type inference
	x := 100
	y := 12.5
	// Type conversion
	var z float64 = float64(x)
	var trung string = strconv.Itoa(x)
	newx, _ := strconv.Atoi(trung)
	name := "Cigarettes"

	fmt.Printf("%d has type is %T\n", x, x)
	fmt.Printf("%f has type is %T\n", y, y)
	fmt.Printf("%f has type is %T\n", z, z)
	fmt.Printf("%s has type is %T\n", trung, trung)
	fmt.Printf("%s has type is %T\n", name, name)
	fmt.Printf("%d has type is %T\n", newx, newx)

	// Zero value
	var age int
	var name11 string
	var active bool

	fmt.Println("Default value age ", age)
	fmt.Println("Default value name ", name11)
	fmt.Println("Default value active ", active)

	// Custom types
	// Type definition : Create a new type
	type Custom int
	var hung Custom = 100
	fmt.Printf("Hung %d has type is %T\n", hung, hung)
	// Type alias : Rename a type
	type RenameInt = int
	var renameInt RenameInt = 100
	fmt.Printf("RenameInt %d has type is %T\n", renameInt, renameInt)
	// Underlying type là kiểu gốc bên dưới của một custom type.
}
