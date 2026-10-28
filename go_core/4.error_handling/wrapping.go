package main

import (
	"errors"
	"fmt"
)

// ErrUserNotFound — sentinel dùng cho demo errors.Is (chi tiết: sentinel.go).
var ErrUserNotFound = errors.New("user not found")

// ValidationError — custom error có trường dữ liệu, dùng cho demo errors.As.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation %s: %s", e.Field, e.Message)
}

// ---- Chuỗi wrap 3 tầng: repo → service → handler ----

func repoFindUser(id int) error {
	// Tầng gốc: tạo error rồi trả lên.
	return ErrUserNotFound
}

func serviceGetUser(id int) error {
	if err := repoFindUser(id); err != nil {
		// %w để BỌC lỗi gốc, bổ sung ngữ cảnh tầng service.
		return fmt.Errorf("service.GetUser(%d): %w", id, err)
	}
	return nil
}

func handlerGetUser(id int) error {
	if err := serviceGetUser(id); err != nil {
		// Wrap thêm lần nữa ở tầng handler.
		return fmt.Errorf("handler: %w", err)
	}
	return nil
}

// ---- errors.Join: gộp NHIỀU lỗi độc lập (Go 1.20+) ----

func validateUser(name, email string) error {
	var errs []error
	if name == "" {
		errs = append(errs, &ValidationError{Field: "name", Message: "không được để trống"})
	}
	if email == "" {
		errs = append(errs, &ValidationError{Field: "email", Message: "không được để trống"})
	}
	// Join trả về một error duy nhất chứa tất cả; nil nếu mọi input đều nil.
	return errors.Join(errs...)
}

func Wrapping() {
	// ---- 1. errors.Unwrap: lột MỘT lớp bọc ----
	wrapped := fmt.Errorf("tang 1: %w", fmt.Errorf("tang 2: %w", ErrUserNotFound))
	fmt.Println("wrapped      :", wrapped)
	fmt.Println("Unwrap(wrapped):", errors.Unwrap(wrapped)) // chỉ lột 1 lớp

	// ---- 2. errors.Is: duyệt TOÀN bộ chuỗi, tìm target ----
	err := handlerGetUser(42)
	fmt.Println("err chain    :", err)
	fmt.Println("Is(ErrUserNotFound):", errors.Is(err, ErrUserNotFound))

	// So sánh sai cách: so sánh chuỗi / == trực tiếp
	fmt.Println("err == ErrUserNotFound:", err == ErrUserNotFound) // false — đã bị wrap

	// ---- 3. errors.As: tìm error theo KIỂU, gán ra biến typed ----
	verr := fmt.Errorf("validate: %w", &ValidationError{Field: "age", Message: "phải >= 0"})
	var vErr *ValidationError
	if errors.As(verr, &vErr) {
		fmt.Printf("As: field=%s, message=%s\n", vErr.Field, vErr.Message)
	}

	// ---- 4. errors.Join: gộp nhiều lỗi validate ----
	if err := validateUser("", ""); err != nil {
		fmt.Println("Join gộp 2 lỗi:")
		// joined error hỗ trợ Unwrap() []error → lấy slice lỗi con rồi range
		if joined, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range joined.Unwrap() {
				fmt.Println("  -", e)
			}
		}
	}
}
