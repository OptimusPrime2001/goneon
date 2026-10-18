# The `error` Interface trong Go

`error` là interface **predeclared** nhỏ nhất của Go — chỉ gồm **một method**.
Bản chất của mọi lỗi trong Go chính là: **một giá trị implement method `Error() string`**.

Khớp node **Error Handling → errors** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/error_basic.go`

> Đọc trước: [Interfaces](<../Methods and Interfaces/Interfaces.md>) —
> bài này dựa hoàn toàn vào concept interface value = (concrete type, value).

---

## 1. Khái niệm & Bản chất

`error` được khai báo ngay trong package `builtin` của Go:

```go
// Source: $GOROOT/src/builtin/builtin.go
type error interface {
    Error() string
}
```

Chỉ 1 method → **bất kỳ type nào có method `Error() string`** đều tự động
implement `error` (implicit implementation — không cần keyword `implements`):

```text
  error interface                      Các "nghi can" implement nó
  ──────────────────────               ─────────────────────────────────────
  Error() string                 ←     *errors.errorString   (errors.New)
                                   ←     *fmt.wrapError       (fmt.Errorf %w)
                                   ←     MyError / *MyError   (custom của bạn)
                                   ←     os.PathError, net.OpError, ...
```

> **Bản chất:** "lỗi" không phải một kiểu dữ liệu đặc biệt — nó chỉ là
> **một interface value** chứa cặp (concrete type, giá trị).
> Kiểu concrete có thể là struct pointer, struct value, thậm chí... `nil`.

---

## 2. Custom Error — struct tự định nghĩa

### Struct value receiver

```go
type MyError struct {
    Op   string
    Code int
}

func (e MyError) Error() string {
    return fmt.Sprintf("%s failed (code %d)", e.Op, e.Code)
}

err := MyError{Op: "withdraw", Code: 402}
fmt.Println(err.Error())
```

```text
withdraw failed (code 402)
```

### Struct pointer receiver (khuyên dùng)

```go
func (e *MyError) Error() string {
    return fmt.Sprintf("%s failed (code %d)", e.Op, e.Code)
}
```

```go
func Withdraw(balance, amount int) (int, error) {
    if amount > balance {
        return balance, &MyError{Op: "withdraw", Code: 402}
    }
    return balance - amount, nil
}

if _, err := Withdraw(100, 300); err != nil {
    fmt.Println("Lỗi:", err)
}
```

Output thật:

```text
Lỗi: withdraw failed (code 402)
```

> **Vì sao dùng pointer receiver?**
> ① Nhất quán với cách custom error phổ biến (stdlib: `*os.PathError`, `*json.SyntaxError`).
> ② Có **identity** — caller so sánh bằng `errors.As` với `*MyError` rõ ràng.
> ③ Tránh copy value lớn mỗi lần gọi `Error()`.
> ⚠️ Nhưng pointer receiver kéo theo **bẫy typed-nil** (mục 4) — phải hiểu rõ.

###Mang theo field có cấu trúc — mang dữ liệu kèm theo

Custom error **mang theo dữ liệu** thay vì chỉ một chuỗi — caller trích xuất được:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation %s: %s", e.Field, e.Message)
}

err := error(&ValidationError{Field: "age", Message: "phải >= 0"})
fmt.Println(err)
```

```text
validation age: phải >= 0
```

> Mang theo structured data là lý do ta hay dùng struct thay cho `errors.New`.
> Cách lấy field ra khỏi error → `errors.As` (bài [Error Wrapping](<Error Wrapping.md>)).

---

## 3. Tại sao `error` chỉ cần 1 method?

So sánh với exception của các ngôn ngữ khác:

```text
  Ngôn ngữ exception                     Go
  ────────────────────────────────────   ────────────────────────────────
  Exception mang theo stack trace, type,       error chỉ mang theo MESSAGE (chuỗi).
  cause... (nặng, nhiều thông tin)       Muốn thêm info → custom struct do BẠN định nghĩa.
  Framework quyết định control flow.     BẠN viết if err != nil → control flow rõ ràng.
```

Sự đơn giản này là **đặc quyền**: interface nhỏ có nghĩa là "mọi thứ đều là error",
dễ compose, dễ wrap, dễ kiểm tra bằng `errors.Is`/`errors.As`.

---

## 4. ⚠️ Bẫy kinh điển: Nil Interface ≠ Interface holding nil

Đây là lỗi **kinh điển nhất** khi làm việc với error trong Go.

```go
func NilTrap(returnNilPointer bool) error {
    if returnNilPointer {
        var e *MyError     // con trỏ nil, nhưng CÓ kiểu (*MyError)
        return e           // interface = (type=*MyError, value=nil) → KHÔNG nil!
    }
    return nil             // interface = (nil, nil) → nil thực sự
}
```

```go
err := NilTrap(true)
fmt.Println("err != nil ?", err != nil)   // sai: caller tưởng CÓ lỗi
fmt.Println("err =", err)                 // in ra "<nil>" → gây nhầm lẫn tột độ
```

Output thật:

```text
err != nil ? true
err = <nil>
```

### Nguyên nhân gốc

```text
  var e *MyError = nil        return nil
         │                          │
         ▼                          ▼
  error interface              error interface
  ┌────────────────────┐       ┌────────────────────┐
  │ type:   *MyError   │       │ type:   nil        │
  │ value:  nil        │       │ value:  nil        │
  └────────────────────┘       └────────────────────┘
  err != nil  →  TRUE ❌        err != nil  →  FALSE ✓
  (interface CÓ type)          (interface thực sự nil)
```

**`err != nil` chỉ kiểm tra interface có *type* nil hay không** —
nếu interface chứa `(type=*MyError, value=nil)` thì `err != nil` vẫn **true**.

### Cách phòng tránh

```go
// ❌ SAI — luôn trả typed pointer, kể cả khi "không có lỗi"
func FindBad(id int) *MyError {
    u, err := find(id)
    if err != nil {
        return &MyError{...}
    }
    return u // u là *MyError nil → caller nhận interface KHÔNG nil
}

// ĐÚNG — luôn trả về error interface, nil đúng nghĩa
func FindGood(id int) error {
    u, err := find(id)
    if err != nil {
        return &MyError{...}
    }
    _ = u
    return nil // error interface (nil, nil) → nil thực sự
}
```

> **Quy tắc vàng:** Hàm trả kiểu **`error`** (interface), **không** trả
> `*MyError` trực tiếp. Khi "không có lỗi", luôn `return nil` cho biến
> có kiểu `error` tường minh.

---

## 5. Best Practices / Lưu ý quan trọng

* `error` chỉ có method `Error() string` — mọi struct có method này **tự là error**.
* **Dùng pointer receiver** `(*T).Error()` cho custom error (theo convention stdlib).
* Mang theo **field có ý nghĩa** (`Op`, `Code`, `Field`...) thay vì chỉ chuỗi — caller
  trích xuất được bằng `errors.As`.
* Hàm trả lỗi **luôn khai báo kiểu trả về `error`** (interface), trả `&T{}` khi
  có lỗi, `nil` đúng nghĩa khi thành công.
* **Tránh typed-nil:** không trả `*T` nil qua interface — kiểm tra `err != nil`
  sẽ "sai" (true) dù thực chất không có lỗi.
* Muốn error hỗ trợ chain/wrap → implement method `Unwrap() error`
  (chi tiết ở [Error Wrapping](<Error Wrapping.md>)).

---

## 6. Summary

* `error` là predeclared interface gồm **1 method** `Error() string` —
  mọi type có method này **tự động implement** (implicit).
* Custom error thường là **struct pointer** (`*MyError`) mang theo **field dữ liệu**
  có ý nghĩa, tốt hơn `errors.New` cho lỗi phức tạp.
* **Nil interface ≠ interface holding nil:** `err != nil` chỉ kiểm tra *type* của
  interface — `(type=*T, value=nil)` vẫn cho `err != nil` = **true** (bẫy kinh điển).
* **Phòng tránh typed-nil:** hàm trả kiểu `error`, thành công → `return nil` tường minh.
* Muốn error tham gia chain wrap → implement `Unwrap() error`.
