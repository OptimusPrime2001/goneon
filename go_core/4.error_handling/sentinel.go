package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ---- Sentinel errors: biến error "danh nghĩa" export sẵn ----
// Quy ước: khai báo package-level, prefix Err, dùng errors.New nguyên sơ.
var (
	ErrNotFound  = errors.New("item not found")
	ErrForbidden = errors.New("access denied")
)

// FindItem trả về sentinel khi không tìm thấy — caller phân nhánh bằng errors.Is.
func FindItem(id int) (string, error) {
	if id != 1 {
		return "", ErrNotFound
	}
	return "apple", nil
}

func Sentinel() {
	// ---- 1. errors.Is với sentinel đã wrap ----
	wrapped := fmt.Errorf("cache lookup: %w", ErrNotFound)
	fmt.Println("Is(wrapped, ErrNotFound):", errors.Is(wrapped, ErrNotFound))
	fmt.Println("== trực tiếp           :", wrapped == ErrNotFound) // false — đã wrap

	// ---- 2. Phân nhánh theo sentinel (control flow bình thường) ----
	if _, err := FindItem(99); err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			fmt.Println("Xử lý riêng: tạo mặc định hoặc 404")
		case errors.Is(err, ErrForbidden):
			fmt.Println("Xử lý riêng: 403")
		default:
			fmt.Println("Lỗi không xác định:", err)
		}
	}

	// ---- 3. Sentinel có sẵn của stdlib ----
	fmt.Println("io.EOF       :", io.EOF)
	fmt.Println("os.ErrNotExist:", os.ErrNotExist)
	// Kiểm tra file không tồn tại bằng errors.Is
	if _, err := os.Open("khong_ton_tai.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Is(err, os.ErrNotExist): true")
	}
}
