# Strings

`string` là kiểu chuỗi trong Go. Chuỗi là `immutable` — không thể sửa trực tiếp từng ký tự trong chuỗi gốc.

- **Bản chất:** Chuỗi là tập hợp các byte **bất biến** xếp liền nhau, mã hóa theo chuẩn **UTF-8**.
- **Bất biến nghĩa là gì?** Khi bạn đã tạo `s := "Hello"`, bạn **không thể** sửa ký tự đầu thành `'X'` bằng `s[0] = 'X'`. Muốn sửa, phải cắt chuỗi và tạo vùng nhớ mới.

```go
package main

import "fmt"

func main() {
	s := "Hello"
	// s[0] = 'X' // ❌ lỗi compile: cannot assign to s[0]

	fmt.Println(s)
}
```

***

## 1. Interpreted string literals (Dấu nháy kép `"..."`)

Chuỗi **được Go diễn dịch (interpreted)** — các ký tự escape như `\n`, `\t` được xử lý thành ký tự thật.

| Escape | Ý nghĩa              |
|--------|-----------------------|
| `\n`   | Xuống dòng            |
| `\t`   | Tab                   |
| `\\`   | Dấu gạch chéo `\`     |
| `\"`   | Dấu nháy kép `"`      |
| `\uXXXX` | Ký tự Unicode (hex) |

```go
package main

import "fmt"

func main() {
	msg := "Xin chào,\ntôi là Dev Go."
	fmt.Println(msg)

	tab := "Name:\tGoneon"
	fmt.Println(tab)

	path := "C:\\Users\\goneon"
	fmt.Println(path)
}
```

```text
Xin chào,
tôi là Dev Go.
Name:	Goneon
C:\Users\goneon
```

Đặc điểm:

- Phải viết trên **1 dòng** (hoặc dùng `\n` để xuống dòng).
- Hỗ trợ escape sequence.
- Phù hợp cho chuỗi thông thường, prompt, format đầu ra.

***

## 2. Raw string literals (Backtick `` `...` ``)

Chuỗi **thô (raw)** — viết sao hiển thị y vậy, Go **không xử lý** escape sequence nào.

```go
package main

import "fmt"

func main() {
	query := `SELECT id, name
FROM users
WHERE age > 18
ORDER BY name;`

	fmt.Println(query)

	jsonData := `{"name": "Goneon", "age": 25}`
	fmt.Println(jsonData)

	path := `C:\Users\goneon` // không cần escape \\
	fmt.Println(path)
}
```

```text
SELECT id, name
FROM users
WHERE age > 18
ORDER BY name;
{"name": "Goneon", "age": 25}
C:\Users\goneon
```

Đặc điểm:

- Hỗ trợ **nhiều dòng** (multiline).
- **Không** xử lý `\n`, `\t`, `\\` — hiển thị nguyên văn.
- Rất hay dùng để viết câu lệnh **SQL**, **JSON thô**, **regex**, **template** trong code.

***

## 3. So sánh 2 dạng chuỗi

| Tiêu chí              | Interpreted `"..."`         | Raw `` `...` ``              |
|------------------------|-----------------------------|------------------------------|
| Dấu ký tự              | Nháy kép `"`                | Backtick `` ` ``             |
| Nhiều dòng             | Cần `\n`                    | ✅ Tự nhiên                  |
| Escape sequence        | ✅ Có (`\n`, `\t`...)       | ❌ Không                     |
| Chứa nháy kép `"`      | Cần escape `\"`             | ✅ Không cần                 |
| Chứa backtick `` ` ``  | ✅ Bình thường              | ❌ Không thể                 |
| Phù hợp               | Chuỗi thông thường          | SQL, JSON, regex, template   |

```go
package main

import "fmt"

func main() {
	a := "Line1\nLine2"   // \n thành xuống dòng thật
	b := `Line1\nLine2`   // \n chỉ là 2 ký tự \ và n

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(len(a)) // 12
	fmt.Println(len(b)) // 12 — cùng độ dài, nhưng a đã xuống dòng
}
```

```text
Line1
Line2
Line1\nLine2
12
12
```

***

## 4. Type conversion giữa string và số

Go **không tự động chuyển** giữa `string` và số. Dùng package `strconv`.

### Số → string

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	n := 42

	s1 := strconv.Itoa(n)                  // int → string
	s2 := strconv.FormatFloat(3.14, 'f', 2, 64) // float64 → string
	s3 := strconv.FormatBool(true)         // bool → string

	fmt.Println(s1) // 42
	fmt.Println(s2) // 3.14
	fmt.Println(s3) // true
}
```

### String → số

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	s := "42"

	n, err := strconv.Atoi(s) // string → int
	if err != nil {
		fmt.Println("lỗi parse:", err)
		return
	}
	fmt.Printf("%T %d\n", n, n) // int 42

	f, err := strconv.ParseFloat("3.14", 64) // string → float64
	if err != nil {
		fmt.Println("lỗi parse:", err)
		return
	}
	fmt.Printf("%T %f\n", f, f) // float64 3.140000
}
```

### Phân biệt conversion và parsing

- **Conversion**: ép kiểu trực tiếp giữa các kiểu số, ví dụ `float64(a)`.
- **Parsing**: chuyển **từ chuỗi sang** kiểu khác, ví dụ `strconv.Atoi("42")`.
- `strconv.Atoi` → parse **string → int**.
- `strconv.Itoa` → parse **int → string**.

### Cảnh báo

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// n, _ := strconv.Atoi("abc") // ❌ lỗi parse: invalid syntax
	n, err := strconv.Atoi("abc")
	fmt.Println(n, err) // 0 strconv.Atoi: parsing "abc": invalid syntax
}
```

Luôn kiểm tra `err` khi parse chuỗi — dữ liệu đầu vào có thể không hợp lệ.

***

## Tóm tắt nhanh

- `string` là **immutable**, mã hóa UTF-8.
- `"..."` → interpreted, xử lý escape, dùng cho chuỗi thường.
- `` `...` `` → raw, không escape, hỗ trợ nhiều dòng, dùng cho SQL/JSON.
- Không tự đổi `string ↔ number` — dùng `strconv`.
- `strconv.Atoi` (string→int), `strconv.Itoa` (int→string).
- Luôn kiểm tra `err` khi parse.

***

## Exercises

### 1. Chọn dạng chuỗi phù hợp

Viết lại mỗi đoạn sau bằng dạng chuỗi phù hợp nhất:

- Một đường dẫn Windows: `C:\Users\goneon\Documents`
- Câu lệnh SQL 3 dòng
- Chuỗi có nội dung: `He said "hello"`

### 2. Word counter

Viết hàm đếm số từ trong một chuỗi (nhiều dòng, dùng raw string):

```go
func wordCount(s string) int
```

Gợi ý: dùng `strings.Fields`.

### 3. Parse và validate

Nhận vào slice `[]string` gồm các số dạng chuỗi, parse thành `int` và in ra những chuỗi **không hợp lệ** kèm thông báo lỗi.
