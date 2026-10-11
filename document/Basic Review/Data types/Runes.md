# Runes

`rune` trong Go là **alias của `int32`**, dùng để đại diện cho **một ký tự Unicode**.

- Một `byte` chỉ chứa được 1 ký tự ASCII.
- Một `rune` chứa được **bất kỳ ký tự nào** trên thế giới (tiếng Việt có dấu, tiếng Nhật, emoji...).

```go
package main

import "fmt"

func main() {
	var c rune = 'A'

	fmt.Println(c)         // 65
	fmt.Printf("%T\n", c)  // int32
	fmt.Println(string(c)) // A
}
```

```text
65
int32
A
```

***

## 1. Rune literal

Rune được viết trong **nháy đơn** `'...'`:

| Kiểu rune literal | Ví dụ        | Ký tự                 |
|-------------------|--------------|------------------------|
| Ký tự thường      | `'A'`        | A                      |
| Escape            | `'\n'`       | Xuống dòng             |
| Hex               | `'\x41'`     | A                      |
| Unicode (4 hex)   | `'\u00e9'`   | é                      |
| Unicode dài       | `'\u00e1'`   | á                      |

```go
package main

import "fmt"

func main() {
	a := 'A'
	e := '\u00e9' // é
	aAcute := '\u00e1' // á
	emoji := '\U0001F600' // 😀

	fmt.Println(a, e, aAcute, emoji)
	fmt.Println(string(a), string(e), string(aAcute), string(emoji))
	fmt.Printf("%d %d %d %d\n", a, e, aAcute, emoji)
}
```

```text
65 233 225 128512
A é á 😀
65 233 225 128512
```

Lưu ý:

- `'A'` và `"A"` **khác nhau**:
  - `'A'` → `rune` (kiểu số, giá trị 65)
  - `"A"` → `string` (kiểu chuỗi, độ dài 1)

***

## 2. byte vs rune

| Tiêu chí         | `byte` (= `uint8`)      | `rune` (= `int32`)           |
|-------------------|-------------------------|------------------------------|
| Kích thước        | 1 byte                  | 4 bytes                      |
| Phạm vi          | 0 → 255                 | -2.147.483.648 → 2.147.483.647 |
| Đại diện          | 1 byte dữ liệu          | 1 ký tự Unicode              |
| ASCII             | ✅ Đủ                    | ✅ Đủ                         |
| Tiếng Việt có dấu | ❌ Không đủ (cần nhiều byte) | ✅ Đủ                  |
| Emoji             | ❌ Không đủ              | ✅ Đủ                         |

```go
package main

import "fmt"

func main() {
	b := byte(65)   // ASCII 'A'
	r := rune('A')  // Unicode 'A'

	fmt.Println(b, r)       // 65 65
	fmt.Printf("%T %T\n", b, r) // uint8 int32
}
```

```text
65 65
uint8 int32
```

***

## 3. UTF-8 — cách Go lưu ký tự

Go dùng chuẩn **UTF-8** để lưu chuỗi. UTF-8 là mã hóa **biến đổi độ dài**:

| Ký tự                  | Số byte | Ví dụ              |
|------------------------|---------|---------------------|
| ASCII (`A`-`Z`, `0`-`9`) | 1 byte | `'A'`, `'1'`      |
| Ký tự Latin mở rộng     | 2 bytes | `'é'`, `'á'`, `'ñ'` |
| Ký tự CJK, Việt có dấu  | 2-3 bytes | `'ế'`, `'ữ'`      |
| Emoji                   | 4 bytes | `'😀'`, `'🚀'`      |

```go
package main

import "fmt"

func main() {
	a := "A"    // 1 byte
	e := "é"    // 2 bytes
	ế := "ế"  // 3 bytes
	emoji := "😀" // 4 bytes

	fmt.Println(len(a))    // 1
	fmt.Println(len(e))    // 2
	fmt.Println(len(ế))    // 3
	fmt.Println(len(emoji)) // 4
}
```

```text
1
2
3
4
```

> **Quan trọng:** `len(s)` trả về **số byte**, KHÔNG phải số ký tự!

***

## 4. Duyệt string bằng range

Khi `range` qua một string, Go trả về **`rune`** và **vị trí byte**:

```go
package main

import "fmt"

func main() {
	s := "Việt Nam"

	for i, r := range s { // i = index byte, r = rune
		fmt.Printf("%d %c %d\n", i, r, r)
	}
}
```

```text
0 V 86
1 i 105
2 ê 234
4 t 116
5   32
7 N 78
8 a 97
9 m 109
```

Lưu ý:

- `i` nhảy cóc (0, 1, 2, **4**, 5, **7**...) vì `'ê'` tốn 2 byte.
- `r` luôn là `rune` hợp lệ — đây là cách **duyệt string an toàn**.

### Cách sai: duyệt bằng index

```go
package main

import "fmt"

func main() {
	s := "Việt Nam"

	// Cách SAI — index trả về byte, không phải ký tự
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c ", s[i]) // s[i] là byte
	}
	fmt.Println()
}
```

```text
V i ?t  N a m
```

Ký tự `'ê'` bị tách thành 2 byte → hiển thị sai.

***

## 5. Đếm số ký tự

Vì `len(s)` trả về số byte, muốn đếm **số ký tự thật** phải dùng `utf8.RuneCountInString` hoặc `range`:

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Xin chào Việt Nam"

	fmt.Println(len(s))                        // số BYTE
	fmt.Println(utf8.RuneCountInString(s))     // số KÝ TỰ

	// Đếm bằng range
	count := 0
	for range s {
		count++
	}
	fmt.Println(count)
}
```

```text
19
17
17
```

***

## 6. Chuyển đổi giữa byte, rune và string

```go
package main

import "fmt"

func main() {
	// rune → string
	r := '😀'
	s := string(r)
	fmt.Println(s) // 😀

	// string → rune (lấy ký tự đầu)
	str := "Hello"
	first := rune(str[0]) // 'H' (vì ASCII)
	fmt.Printf("%c %d\n", first, first)

	// string → []rune (tất cả ký tự)
	runes := []rune("Việt")
	fmt.Println(runes) // [86 105 234 116]
	fmt.Println(string(runes)) // Việt

	// byte → string
	b := byte(72)
	fmt.Println(string(b)) // H
}
```

```text
😀
H 72
[86 105 234 116]
Việt
H
```

***

## Tóm tắt nhanh

- `rune` = `int32` — đại diện 1 ký tự Unicode.
- `byte` = `uint8` — đại diện 1 byte dữ liệu.
- Rune literal viết trong **nháy đơn**: `'A'`, `'\u00e9'`.
- UTF-8 có độ dài biến đổi: ASCII 1 byte, tiếng Việt 2-3 byte, emoji 4 byte.
- `len(s)` trả về **số byte**, không phải số ký tự.
- Duyệt string an toàn bằng `range` → nhận `rune`.
- Đếm ký tự: `utf8.RuneCountInString(s)` hoặc `for range s`.

***

## Exercises

### 1. byte hay rune?

Chỉ ra kiểu của mỗi biểu thức sau và giá trị của nó:

```go
'A'
"😀"
'\n'
byte(65)
rune('B')
```

### 2. Đếm ký tự tiếng Việt

Viết hàm đếm số ký tự (không phải byte) của một chuỗi tiếng Việt:

```go
func countRunes(s string) int
```

Yêu cầu:

- Không dùng `len()`.
- Chạy thử với `"Học lập trình Go"`.

### 3. Đảo chuỗi an toàn

Viết hàm đảo ngược chuỗi, đảm bảo các ký tự Unicode (tiếng Việt, emoji) **không bị hỏng**:

```go
func reverse(s string) string
```

Gợi ý: chuyển sang `[]rune` trước khi đảo.

Ví dụ:

- `"Hello"` → `"olleH"`
- `"Việt Nam"` → `"maN tiệV"`
- `"Go🚀"` → `"🚀oG"`
