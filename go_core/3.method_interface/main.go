package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	id, _ := uuid.NewRandom()
	fmt.Println("Method and Interface\n", id)
	Methods()
	Interfaces()
	Generics()
}
