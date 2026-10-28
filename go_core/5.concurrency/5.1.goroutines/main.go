package main

import (
	"fmt"
	"time"
)

func printNumber() {
	for i := 1; i <= 5; i++ {
		fmt.Println("In số: ", i)
	}
}
func printLetter() {
	for char := 'A'; char <= 'E'; char++ {
		fmt.Printf("%c\n", char)
	}
}
func BAI1() {
	fmt.Println("---- BAI1 ----")
	go printNumber()
	printLetter()
}
func sayHello() {
	for i := range 3 {
		fmt.Println("Hello", i)
	}
}
func BAI2() {
	fmt.Println("---- BAI2 ----")
	go sayHello()
	// sayHello()
	fmt.Println("Main finished")
}
func printMessage(mess string, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("Message %s lần thứ %d\n", mess, i)
	}
}

func BAI3() {
	fmt.Println("---- BAI3 ----")
	go printMessage("Hello Go", 3)
	go printMessage("Learning Goroutine", 2)
}
func BAI4() {
	go func() {
		fmt.Println("Task started")
		time.Sleep(time.Second * 2)
		fmt.Println("Task completed")
	}()
}
func BAI5() {
	for id := 1; id < 6; id++ {
		go func(id int) {
			fmt.Printf("Worker %d started\n", id)
			for count := 1; count < 4; count++ {
				fmt.Printf("Worker %d - interation %d \n", id, count)
			}
			fmt.Printf("Worker %d finished \n", id)
		}(id)
	}
}
func main() {
	// BAI1()
	// BAI2()
	// BAI3()
	// BAI4()
	// BAI5()
	fmt.Println("Main finished")
	time.Sleep(time.Millisecond)
}
