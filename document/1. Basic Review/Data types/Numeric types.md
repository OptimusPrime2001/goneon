# Numeric Types

Go hỗ trợ nhiều kiểu số, phân chia rất tỉ mỉ dựa trên **kích thước bộ nhớ (bit)** và **dấu (+/-)** để lập trình viên tối ưu hóa hiệu năng.

| Nhóm          | Kiểu                                              |
|---------------|---------------------------------------------------|
| Số nguyên có dấu     | `int8`, `int16`, `int32`, `int64`, `int`   |
| Số nguyên không dấu  | `uint8`, `uint16`, `uint32`, `uint64`, `uint` |
| Số thực              | `float32`, `float64`                          |
| Số phức              | `complex64`, `complex128`                     |

***

## 1. Integers (Số nguyên)

### Signed Integers — số nguyên có dấu

Có thể mang giá trị **âm hoặc dương**:

| Kiểu     | Kích thước | Phạm vi                    |
|----------|-----------|-----------------------------|
| `int8`   | 8-bit (1 byte)   | -128 → 127            |
| `int16`  | 16-bit (2 bytes) | -32.768 → 32.767      |
| `int32`  | 32-bit (4 bytes) | -2.147.483.648 → 2.147.483.647 |
| `int64`  | 64-bit (8 bytes) | -9.223.372.036.854.775.808 → 9.223.372.036.854.775.807 |
| `int`    | 32 hoặc 64-bit (theo OS) | phụ thuộc nền tảng |

```go
package main

import "fmt"

func main() {
	var a int8 = 127
	var b int32 = 2147483647
	var c int64 = 9223372036854775807

	fmt.Println(a, b, c)
	fmt.Printf("%d %d %d\n", a, b, c)
}
```

```text
127 2147483647 9223372036854775807
```

### Unsigned Integers — số nguyên không dấu

Chỉ mang giá trị **từ 0 trở lên**, do đó phạm vi dương lớn gấp đôi so với kiểu có dấu cùng kích thước:

| Kiểu      | Kích thước | Phạm vi                |
|-----------|-----------|-------------------------|
| `uint8`   | 8-bit     | 0 → 255                 |
| `uint16`  | 16-bit    | 0 → 65.535              |
| `uint32`  | 32-bit    | 0 → 4.294.967.295       |
| `uint64`  | 64-bit    | 0 → 18.446.744.073.709.551.615 |
| `uint`    | 32 hoặc 64-bit | phụ thuộc nền tảng |

```go
package main

import "fmt"

func main() {
	var a uint8 = 255
	var b uint16 = 65535

	fmt.Println(a, b)

	// a++ // ❌ tràn số: 255 + 1 = 0 (overflow)
}
```

### Kiểu số nguyên đặc biệt

#### `int` và `uint`

Kích thước **không cố định**:

- OS 64-bit → `int` = `int64`
- OS 32-bit → `int` = `int32`

> **Thực tế:** 90% trường hợp viết code thông thường (vòng lặp, đếm số, ID...), dev Go chỉ dùng `int`.

#### `byte`

Thực chất là **alias** của `uint8`. Khi đọc file, truyền dữ liệu qua mạng, bạn đang làm việc với mảng `byte`.

```go
package main

import "fmt"

func main() {
	var b byte = 65
	fmt.Println(b)      // 65
	fmt.Println(string(b)) // A
}
```

#### Alias khác

- `rune` = `int32` → xem chi tiết tại [Runes.md](./Runes.md)
- `uintptr` → kiểu con trỏ (hiếm dùng, chỉ trong thư viện hệ thống)

***

## 2. Floating Point (Số thực)

Dùng cho số có dấu phẩy thập phân. Go **không có** kiểu `float` chung chung mà bắt buộc chọn kích thước:

| Kiểu       | Độ chính xác     | Ghi chú                              |
|------------|------------------|---------------------------------------|
| `float32`  | Single precision | Ít dùng, dễ sai số với số thập phân dài |
| `float64`  | Double precision | **Nên dùng mặc định**                 |

```go
package main

import "fmt"

func main() {
	var a float32 = 3.14159265358979
	var b float64 = 3.14159265358979

	fmt.Printf("%.10f\n", a)
	fmt.Printf("%.10f\n", b)
}
```

```text
3.1415927410
3.1415926536
```

> **Thực tế:** Luôn dùng `float64` khi làm việc với số thực trong Go, trừ khi hệ thống cực hạn cần tiết kiệm từng chút RAM.

### Cạm bẫy so sánh số thực

```go
package main

import "fmt"

func main() {
	x := 0.1
	y := 0.2
	fmt.Println(x + y == 0.3) // ❌ false — sai số tròn lẻ
	fmt.Println(x+y-0.3 < 1e-9) // ✅ true — so sánh với epsilon
}
```

```text
false
true
```

***

## 3. Complex Numbers (Số phức)

Go có `complex64` và `complex128` cho số phức. Cực kỳ hiếm gặp trong lập trình ứng dụng thông thường (chủ yếu dùng trong khoa học, xử lý tín hiệu).

| Kiểu         | Thành phần thực / ảo |
|--------------|----------------------|
| `complex64`  | `float32`            |
| `complex128` | `float64`            |

```go
package main

import "fmt"

func main() {
	var c complex128 = complex(3, 4) // 3 + 4i

	fmt.Println(c)
	fmt.Println(real(c)) // 3
	fmt.Println(imag(c)) // 4
	fmt.Println(c * c)   // -7 + 24i
}
```

```text
(3+4i)
3
4
(-7+24i)
```

***

## 4. Type Conversion số học

Go **không tự động ép kiểu** số. Bạn phải ép kiểu tường minh.

### int → float64

```go
package main

import "fmt"

func main() {
	var a int = 10
	var b float64 = float64(a)

	fmt.Printf("%T %T\n", a, b)
	fmt.Println(b)
}
```

```text
int float64
10
```

### float64 → int (mất phần thập phân)

```go
package main

import "fmt"

func main() {
	x := 3.99
	y := int(x) // cắt phần thập phân, KHÔNG làm tròn

	fmt.Println(y) // 3
}
```

### Cảnh báo overflow

```go
package main

import "fmt"

func main() {
	var a int8 = 127
	// a++ // ❌ lỗi runtime: overflow (ngoài phạm vi -128 → 127)

	fmt.Println(a)
}
```

Lưu ý:

- `float64(a)` là **conversion** (ép kiểu trực tiếp).
- Ép `float64 → int` **cắt** phần thập phân, không làm tròn.
- Ép kiểu vượt phạm vi → lỗi runtime `overflow`.

***

## Tóm tắt nhanh

- Dùng `int` cho số nguyên thông thường.
- Dùng `uint` khi giá trị **không bao giờ âm** (kích thước, ID, bộ đếm).
- Dùng `float64` cho số thực.
- `byte` = `uint8`, `rune` = `int32`.
- Go **không ép kiểu ngầm định** — phải viết `float64(a)` hay `int(x)` tường minh.
- `complex64/128` chỉ dùng khi thật sự cần số phức.

***

## Exercises

### 1. Kiểm tra phạm vi

Viết chương trình in ra:

- Phạm vi nhỏ nhất và lớn nhất của `int8`, `uint8`, `int16`, `uint16`.
- Dùng `math` package: `math.MinInt8`, `math.MaxInt8`...

### 2. Tràn số có chủ đích

```go
var a uint8 = 255
a++ // giá trị sau khi tăng là bao nhiêu? vì sao?
```

Chạy thử và giải thích kết quả.

### 3. Chuyển đổi kiểu

- Viết hàm nhận `int` và trả về `float64`.
- Viết hàm nhận `float64` và trả về `int` (cắt thập phân).
- In kết quả `3.99 → int` và giải thích vì sao được `3` chứ không phải `4`.
