package main

import (
	"errors"
	"net/mail"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   *int   `json:"age"` // con trỏ: phân biệt "không gửi" và 0
}

func (r *CreateUserRequest) Validate() []FieldError {
	// Chuẩn hóa trước khi kiểm tra
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	var errs []FieldError
	add := func(field, msg string) { errs = append(errs, FieldError{field, msg}) }

	switch n := utf8.RuneCountInString(r.Name); {
	case n == 0:
		add("name", "is required")
	case n < 2 || n > 50:
		add("name", "must be between 2 and 50 characters")
	}

	if r.Email == "" {
		add("email", "is required")
	} else if !isEmail(r.Email) {
		add("email", "must be a valid email")
	}

	if r.Age != nil && (*r.Age < 0 || *r.Age > 150) {
		add("age", "must be between 0 and 150")
	}
	return errs
}

func isEmail(s string) bool {
	a, err := mail.ParseAddress(s)
	return err == nil && a.Address == s // loại bỏ dạng "Tên <a@b.com>"
}

/// Sử dụng thư viện validator

type CreateUserRequest2 struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   *int   `json:"age" validate:"omitempty,gte=0,lte=150"`
}

func validate() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}
func validateStruct(s any) []FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var ves validator.ValidationErrors
	if !errors.As(err, &ves) {
		return []FieldError{{Field: "", Message: "invalid input"}}
	}
	out := make([]FieldError, 0, len(ves))
	for _, fe := range ves {
		out = append(out, FieldError{Field: fe.Field(), Message: messageFor(fe)})
	}
	return out
}
func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "must be at most " + fe.Param()
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	default:
		return "is invalid"
	}
}
