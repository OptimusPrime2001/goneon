# Data Types

Go là ngôn ngữ `static type`, nghĩa là mỗi biến đều có kiểu dữ liệu xác định từ lúc khai báo. Go không ép kiểu ngầm định — bạn phải khai báo hoặc để Go suy luận kiểu một cách tường minh.

Mục tiêu của phần này:

- Nắm được các kiểu dữ liệu cơ bản trong Go.
- Phân biệt số nguyên có dấu / không dấu, chuỗi interpreter / raw.
- Hiểu `rune` và cách Go mã hóa ký tự theo UTF-8.

***

## Các kiểu dữ liệu cơ bản

| Kiểu nhóm        | Kiểu cụ thể                                    | File chi tiết      |
|------------------|-------------------------------------------------|--------------------|
| Numeric          | `int`, `uint`, `int8..64`, `uint8..64`, `byte`  | [Numeric types.md](./Numeric%20types.md) |
|                  | `float32`, `float64`                            | [Numeric types.md](./Numeric%20types.md) |
|                  | `complex64`, `complex128`                       | [Numeric types.md](./Numeric%20types.md) |
| String           | `"..."` (interpreted), `` `...` `` (raw)        | [Strings.md](./Strings.md) |
| Boolean          | `bool`                                          | [Boolean.md](./Boolean.md) |
| Rune             | `rune` (= `int32`)                              | [Runes.md](./Runes.md) |

***

## Mục lục

1. [Numeric types](./Numeric%20types.md) — số nguyên (signed/unsigned), số thực, số phức
2. [Strings](./Strings.md) — interpreted literal, raw literal, chuyển đổi kiểu chuỗi
3. [Boolean](./Boolean.md) — `true` / `false`, không ép kiểu ngầm định
4. [Runes.md](./Runes.md) — ký tự Unicode, `byte` vs `rune`

***

## Ví dụ tổng quan

```go
package main

import "fmt"

func main() {
	age := 25                // int
	price := 19.99           // float64
	name := "Goneon"         // string
	isActive := true         // bool
	grade := 'A'             // rune
	b := byte(72)            // byte (= uint8)
	c := complex(3, 4)       // complex128

	fmt.Printf("%T %T %T %T %T %T %T\n",
		age, price, name, isActive, grade, b, c)
}
```

```text
int float64 string bool rune uint8 complex128
```

***

## Gợi ý cách học

- Đọc từng file theo thứ tự: `Numeric types` → `Strings` → `Boolean` → `Runes`.
- Chạy lại toàn bộ ví dụ bằng `go run`.
- Tự trả lời: kiểu nào dùng cho trường hợp nào, vì sao Go không ép kiểu ngầm định.
