package main

import (
	"errors"
	"fmt"
)

// Divide minh họa idiomatic: error LUÔN là tham số trả về cuối cùng.
// Trả về (zero value, error) khi thất bại; (value, nil) khi thành công.
func Divide(a, b int) (int, error) {
	if b == 0 {
		// errors.New: tạo error nguyên sơ (static message, không %w).
		return 0, errors.New("divide: chia cho 0")
	}
	return a / b, nil
}

// MyError — custom error kiểu struct.
// Method Error() có receiver POINTER → *MyError implement error interface.
type MyError struct {
	Op   string
	Code int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s failed (code %d)", e.Op, e.Code)
}

// Withdraw rút tiền; trả về custom error khi số dư không đủ.
func Withdraw(balance, amount int) (int, error) {
	if amount > balance {
		return balance, &MyError{Op: "withdraw", Code: 402} // 402 = Payment Required
	}
	return balance - amount, nil
}

// NilTrap minh họa "typed nil" — cái bẫy kinh điển của error interface.
func NilTrap(returnNilPointer bool) error {
	if returnNilPointer {
		var e *MyError // con trỏ nil, nhưng CÓ kiểu
		return e       // interface chứa (type=*MyError, value=nil) → KHÔNG phải nil interface
	}
	return nil // interface (nil, nil) → nil thực sự
}

func ErrorBasic() {
	// ---- 1. Pattern chuẩn: gọi → kiểm tra err NGAY, fail fast ----
	if result, err := Divide(10, 2); err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("10 / 2 =", result)
	}

	if _, err := Divide(10, 0); err != nil {
		fmt.Println("Lỗi:", err)
	}

	// ---- 2. Custom error struct ----
	if newBalance, err := Withdraw(100, 300); err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("Số dư còn:", newBalance)
	}

	// ---- 3. Bẫy typed nil: err != nil nhưng "không có lỗi" ----
	err := NilTrap(true)
	fmt.Println("err != nil:", err != nil) // sai: err KHÔNG nil dù value là nil pointer
	fmt.Println("err =", err)              // in ra "<nil>" vì method gọi trên nil receiver
}
