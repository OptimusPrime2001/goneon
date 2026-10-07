package main

import "fmt"

func Boolean() bool {
	fmt.Println("\n--------Boolean types----------")
	type Config struct {
		Debug   bool
		Verbose bool
		Enabled bool
	}
	var isReady = false
	isActive := true
	cfg := Config{
		Debug:   true,
		Verbose: false,
		Enabled: true,
	}
	if cfg.Debug {
		fmt.Println("Chế độ debug đang bật")
	}

	return (isReady && isActive)
}

func Numeric() {
	fmt.Println("\n--------------Numeric types-------------")
	var a int8 = -127
	var b int32 = 124324234
	var c int64 = 324324234234231432
	var d uint8 = 128
	e := 32.223
	var x int = int(e)
	fmt.Println(a, b, c, d, e, x)
	fmt.Printf("%T %T %T \n", a, b, c)
}
func Runes() {
	fmt.Println("\n--------------Runes types-------------")
	a := 'A'
	e := '\u00e9'         // é
	aAcute := '\u00e1'    // á
	emoji := '\U0001F600' // 😀
	s := "Việt Nam"

	for i, r := range s { // i = index byte, r = rune
		fmt.Printf("%d %c %d\n", i, r, r)
	}
	fmt.Println(a, e, aAcute, emoji)
	fmt.Println(string(a), string(e), string(aAcute), string(emoji))
	fmt.Printf("%d %d %d %d\n", a, e, aAcute, emoji)
}
func DataType() {
	bool := Boolean()

	fmt.Println("Check boolean", bool)

	Numeric()
	Runes()
}
