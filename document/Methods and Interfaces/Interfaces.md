# Interfaces trong Go (Golang)

Interface là **tập hợp method** — type nào có **đủ** method là tự động implement (không cần keyword).

Khớp node **Interfaces** của roadmap.sh (§6).
Code mẫu: `fundamental/4.method_interface/interface.go`

---

## 1. Interface Basic

```go
type Shape interface {
    Area() float64
}
```

### Implicit implementation — không có `implements`

```go
type Circle struct{ Radius float64 }

func (c Circle) Area() float64 {
    return 3.14 * c.Radius * c.Radius
}
// Circle đã implement Shape — CHỈ cần có Area() float64, không khai báo gì thêm
```

```go
func PrintArea(s Shape) {         // nhận MỌI type nào có Area() float64
    fmt.Println("Area:", s.Area())
}

PrintArea(Circle{Radius: 2})      // Area: 12.56
```

Output thật:

```text
Area: 12.56
Area: 6
```

> Go đảo ngược so với Java/C#: **type tự khai báo mình implement** → viết interface "không hợp đồng" (không cần biết trước type nào sẽ dùng).

---

## 2. Interface Value = (type, value)

Interface không lưu dữ liệu trực tiếp mà giữ **cặp (concrete type, giá trị)**:

```text
s (Shape)              internal representation
┌──────────────┐
│ type: Circle │──┐
│ value: 2.0   │  │    ┌──────────────────┐
└──────────────┘  └──▶ │ Circle{Radius:2} │
                        └──────────────────┘
```

### Nil interface ≠ interface holding nil

```go
var s Shape                          // s = (nil, nil)        → gọi method = PANIC
var c *Circle = nil
s = c                               // s = (*Circle, nil)    → method nhận nil pointer
```

* Kiểm tra: `if s == nil` chỉ đúng với trạng thái **đầu tiên**
* Trap kinh điển: trả `*T` nil từ hàm → interface "có type" nên **không** nil → bug khó tìm

---

## 3. Empty Interface — `any`

```go
var v any          // viết tắt của interface{} — KHÔNG method nào bắt buộc
v = 42
v = "hello"
v = Circle{}
```

```go
printAnything(42)
printAnything("hello")
```

Output thật:

```text
any: 42
any: hello
```

* `any` chấp nhận **mọi** giá trị — dùng khi chưa biết type (vd: `fmt.Println`, `json.Unmarshal`)
* Muốn dùng giá trị → cần **type assertion** (mục 5)

---

## 4. Embedding Interface — gộp interface

```go
type Reader interface{ Read() string }
type Writer interface{ Write(s string) }

// ReadWriter = Reader + Writer — embedding, không viết lại method
type ReadWriter interface {
    Reader
    Writer
}
```

```go
var rw ReadWriter = File{name: "note.txt"}   // phải có ĐỦ Read + Write
fmt.Println(rw.Read())
rw.Write("hello")
```

Output thật:

```text
data from note.txt
write hello to note.txt
```

> Đây là mô hình của stdlib: `io.ReadWriter` = `io.Reader` + `io.Writer` — **interface nhỏ**, gộp lại thành lớn.

---

## 5. Type Assertion — "unbox" interface

Lấy concrete type ra khỏi interface:

```go
s, ok := v.(string)     // Comma-ok: an toàn, không panic
if ok {
    fmt.Println(s)      // dùng được như string
} else {
    fmt.Println("không phải string")
}
```

```go
typeAssertion("xin chao")
typeAssertion(3.14)
```

Output thật:

```text
assert string: xin chao
assert: không phải string
```

### Không dùng comma-ok → panic nếu sai type

```go
v := 3.14
s := v.(string)         // ❌ panic: interface conversion: interface {} is float64, not string
```

---

## 6. Type Switch — switch theo concrete type

```go
switch v.(type) {
case int:
    fmt.Println("type: int")
case string:
    fmt.Println("type: string")
case Circle:
    fmt.Println("type: Circle")
default:
    fmt.Println("type: unknown")
}
```

Output thật:

```text
type: int
type: string
type: Circle
```

* `v.(type)` **chỉ dùng được trong type switch** — không có ở nơi khác
* Gán giá trị trong case: `switch x := v.(type)` → `x` có kiểu tương ứng từng nhánh

---

## 7. Summary

* Interface = tập hợp method; **implicit implementation** — có đủ method là implement, không keyword `implements`
* Interface value = **(concrete type, value)**; nil interface ≠ interface holding nil pointer (trap kinh điển)
* `any` = `interface{}` = empty interface — chấp nhận mọi giá trị, dùng kèm type assertion
* **Embedding** gộp interface nhỏ thành lớn (mô hình `io.ReadWriter`)
* **Type assertion** `v.(T)`: dùng comma-ok `v, ok := ...` để tránh panic
* **Type switch** `switch v.(type)`: rẽ nhánh theo concrete type
* Go chuộng **interface nhỏ** (1-2 method) — dễ implement, dễ test
