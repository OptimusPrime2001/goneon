# Boolean

Kiểu `bool` trong Go chỉ có **2 giá trị**:

- `true`
- `false`

```go
package main

import "fmt"

func main() {
	isReady := true
	isDone := false

	fmt.Println(isReady, isDone)
	fmt.Printf("%T %T\n", isReady, isDone)
}
```

```text
true false
bool bool
```

***

## 1. Không có ép kiểu ngầm định

Đây là quy tắc **cực kỳ nghiêm ngặt** của Go:

> Ở C hay JavaScript, số `0` là `false`, số `1` là `true`. Trong Go, **không có chuyện đó**. `0` là `int`, không liên quan gì đến `bool`.

```go
package main

import "fmt"

func main() {
	// if 1 { // ❌ lỗi compile: non-boolean condition
	// 	fmt.Println("yes")
	// }

	x := 1
	if x != 0 { // ✅ phải viết điều kiện boolean tường minh
		fmt.Println("x khác 0")
	}
}
```

```text
x khác 0
```

Lưu ý:

- `if x { ... }` với `x` là `int` → **lỗi compile**.
- Không có phép ép kiểu ngầm `int → bool`.
- Phải so sánh tường minh: `x != 0`, `x > 0`, `len(s) > 0`.

***

## 2. Toán tử so sánh

Toán tử so sánh trả về giá trị `bool`:

| Toán tử | Ý nghĩa          |
|---------|------------------|
| `==`    | Bằng             |
| `!=`    | Không bằng       |
| `<`     | Nhỏ hơn          |
| `>`     | Lớn hơn          |
| `<=`    | Nhỏ hơn hoặc bằng |
| `>=`    | Lớn hơn hoặc bằng |

```go
package main

import "fmt"

func main() {
	a := 10
	b := 20

	fmt.Println(a == b) // false
	fmt.Println(a != b) // true
	fmt.Println(a < b)  // true
	fmt.Println(a >= b) // false
}
```

```text
false
true
true
false
```

***

## 3. Toán tử logic

| Toán tử | Ý nghĩa                |
|---------|------------------------|
| `&&`    | AND — cả 2 phải đúng   |
| `\|\|`  | OR — ít nhất 1 đúng    |
| `!`     | NOT — đảo ngược giá trị |

```go
package main

import "fmt"

func main() {
	a := true
	b := false

	fmt.Println(a && b) // false
	fmt.Println(a || b) // true
	fmt.Println(!a)     // false
	fmt.Println(!b)     // true
}
```

```text
false
true
false
true
```

### Short-circuit evaluation

Go **dừng đánh giá** ngay khi kết quả đã xác định:

```go
package main

import "fmt"

func checkA() bool {
	fmt.Println("checkA chạy")
	return true
}

func checkB() bool {
	fmt.Println("checkB chạy")
	return false
}

func main() {
	// false && x → x không bao giờ được gọi
	r1 := false && checkB() // checkB KHÔNG chạy
	fmt.Println(r1)

	// true || x → x không bao giờ được gọi
	r2 := true || checkA() // checkA KHÔNG chạy
	fmt.Println(r2)
}
```

```text
false
true
```

Ý nghĩa:

- Tránh gọi hàm tốn kém khi không cần thiết.
- Rất hay dùng để kiểm tra `nil` trước khi dereference con trỏ.

***

## 4. bool trong điều kiện

### if / else

```go
package main

import "fmt"

func main() {
	age := 20

	if age >= 18 {
		fmt.Println("Đủ tuổi")
	} else {
		fmt.Println("Chưa đủ tuổi")
	}
}
```

```text
Đủ tuổi
```

### for với điều kiện bool

```go
package main

import "fmt"

func main() {
	i := 0

	for i < 3 {
		fmt.Println(i)
		i++
	}
}
```

```text
0
1
2
```

### Trả về bool từ hàm

```go
package main

import "fmt"

func isEven(n int) bool {
	return n%2 == 0
}

func main() {
	fmt.Println(isEven(4)) // true
	fmt.Println(isEven(7)) // false
}
```

```text
true
false
```

***

## 5. bool làm cờ (flag)

Trong thực tế, `bool` thường dùng để biểu thị trạng thái bật/tắt:

```go
package main

import "fmt"

type Config struct {
	Debug    bool
	Verbose  bool
	Enabled  bool
}

func main() {
	cfg := Config{
		Debug:   true,
		Verbose: false,
		Enabled: true,
	}

	if cfg.Debug {
		fmt.Println("Chế độ debug đang bật")
	}
}
```

```text
Chế độ debug đang bật
```

***

## Tóm tắt nhanh

- `bool` chỉ có `true` và `false`.
- Go **không ép kiểu ngầm định** — `if 1 {}` là lỗi compile.
- Toán tử so sánh (`==`, `<`, `>=`...) trả về `bool`.
- Toán tử logic: `&&`, `||`, `!`.
- Go có **short-circuit evaluation** — dừng đánh giá khi kết quả đã rõ.
- Dùng `bool` làm cờ trạng thái trong cấu trúc dữ liệu.

***

## Exercises

### 1. Tìm lỗi compile

Vì sao các đoạn sau bị lỗi? Sửa lại cho hợp lệ:

```go
x := 0
if x {
	fmt.Println("zero")
}
```

```go
age := 15
if age = 18 {
	fmt.Println("adult")
}
```

### 2. Phép tính logic

Không chạy code, viết kết quả:

```go
fmt.Println(true && false || true)
fmt.Println(!true && false)
fmt.Println(false || !false && true)
```

### 3. Validator

Viết hàm kiểm tra một số có hợp lệ là điểm thi hay không:

```go
func isValidScore(score int) bool
```

Yêu cầu:

- Điểm từ 0 đến 100.
- Trả về `true` nếu hợp lệ, `false` nếu không.
- Viết thêm hàm `isPass(score int) bool` — điểm >= 50 thì đạt.
