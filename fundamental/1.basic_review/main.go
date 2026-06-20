package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Learning()
	TempConverter()
	CurrencyConverter()
}
func Learning() {
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

type Celsius float64
type Fahrenheit float64
// Celsius -> Fahrenheit
func cToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}
// Fahrenheit -> Celsius
func fToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}
func TempConverter() {

	c := Celsius(25)
	f := cToF(c)

	fmt.Printf("%.1f°C = %.1f°F\n", c, f)

	f2 := Fahrenheit(77)
	c2 := fToC(f2)

	fmt.Printf("%.1f°F = %.1f°C\n", f2, c2)
}


type USD float64
type VND float64

func (u USD) ToVND() VND {
	return VND(float64(u) * 25000)
}

func (v VND) ToUSD() USD {
	return USD(float64(v) / 25000)
}

func CurrencyConverter() {
	vnd := VND(2500000)
	fmt.Printf("vnd : %v\n", vnd)
	usd := vnd.ToUSD()
	fmt.Printf("%.1fVND = %.1fUSD\n", vnd, usd)
}