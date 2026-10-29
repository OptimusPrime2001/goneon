package main

import "fmt"

type MyError struct {
	Code int
}

func (e *MyError) Error() string {
	return "my error"
}

func test() error {
	return fmt.Errorf("test failed: %w", &MyError{Code: 500})
}
func main() {
	myErr := &MyError{Code: 500}

	// Tạo đối tượng, chưa gọi Error()
	fmt.Printf("Code: %d\n", myErr.Code)

	// Gọi Error() tường minh
	fmt.Println(myErr.Error())

	// fmt.Errorf tự sử dụng Error() để định dạng lỗi
	// err := fmt.Errorf("test failed: %w", myErr)
	fmt.Println(myErr)
}
