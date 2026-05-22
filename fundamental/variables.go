package fundamental

import (
	"fmt"
)

func Variables() {
	// 1: Using var keyword
	// 1.1 Declare a variable with a type and assign a value
	var name string = "Le Van Trung"
	var age int = 25
	var isActive bool = true
	// 1.2 Declare without type inference
	var score = 90
	var company = "Google"
	// 1.3 Declare without initialization (Zero Value)
	var age2 int     // age2 = 0
	var isReady bool // isReady = false
	// 1.4 Declare multiple variables at once
	var a, b, c int
	fmt.Println("Example 1", a, b, c)

	var x, y, z = 1, 2.5, "Hello"
	fmt.Println("Example 2", x, y, z)
	// 1.5 Declare with group block (Block)
	var (
		id        int
		userName  string = "admin"
		isPremium bool   = true
	)
	fmt.Println("Variables", name, id, userName, isPremium, age, isActive, company, score, age2, isReady)
	// 2: Using operator := ( Short Variable Declaration Operator)
	// score := 120
	message := "Good Job"
	length, width := 10, 20
	length, height := 25, 30
	fmt.Println("Logger variable", score, message, length, width, height)

}

/*
>>> Các lưu ý về variable trong go
1: Biến được khai báo bằng := chỉ hợp lệ trong block và function, ko dùng ở chế độ global
2: Từ khoá var ko được khai báo lại trong cùng một block, từ khoá := được phép nếu đi kèm vói ít nhất 1 biến mới
*/
