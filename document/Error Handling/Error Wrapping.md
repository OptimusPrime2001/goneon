# Error Wrapping & Unwrapping: `errors.Is`, `errors.As`

Khi một lỗi đi qua **nhiều tầng** (handler → service → repo), ta **wrap** nó ở mỗi tầng
để thêm ngữ cảnh, tạo nên một **chuỗi lỗi** (error chain). Go cung cấp bộ công cụ để
**bọc** (`%w`) và **lột/giải nén** (`Unwrap`, `Is`, `As`) chuỗi này.

Khớp node **Error Handling → Error wrapping / errors.Is / errors.As** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/wrapping.go`

> Go 1.13+ giới thiệu wrapping; Go 1.20+ thêm nhiều `%w` và `errors.Join`.

---

## 1. Khái niệm & Bản chất

```text
  Tầng repo        Tầng service        Tầng handler         Caller
  ───────────      ──────────────      ──────────────        ──────────
  errors.New()  →  fmt.Errorf %w   →   fmt.Errorf %w    →   errors.Is/As
  (gốc)           (bọc + ngữ cảnh)    (bọc + ngữ cảnh)     (kiểm tra chuỗi)
  ─────────────────────────────────────────────────────────────────────
                     error chain (chuỗi lỗi) → caller vẫn thấy được NGUỒN GỐC
```

**Bản chất:** `%w` tạo ra một error **giữ con trỏ tới error cha**.
Gọi `Unwrap()` sẽ "lột" từng lớp, dần dần về tới error gốc.

```go
// Mỗi tầng chỉ cần: nhận err, bọc ngữ cảnh bằng %w, trả lên.
func repoFindUser(id int) error        { return ErrUserNotFound }
func serviceGetUser(id int) error      { return fmt.Errorf("service.GetUser(%d): %w", id, repoFindUser(id)) }
func handlerGetUser(id int) error      { return fmt.Errorf("handler: %w", serviceGetUser(id)) }
```

```text
handlerGetUser(42) tạo chuỗi:
  "handler: service.GetUser(42): user not found"
        │                    │              │
        │                    │              └── nguồn gốc: ErrUserNotFound
        │                    └── ngữ cảnh tầng service
        └── ngữ cảnh tầng handler
```

---

## 2. `errors.Unwrap` — lột từng lớp

`errors.Unwrap(err)` trả về error **cha trực tiếp** (một lớp), hoặc `nil` nếu không wrap:

```go
wrapped := fmt.Errorf("tầng 1: %w", fmt.Errorf("tầng 2: %w", ErrUserNotFound))
fmt.Println("wrapped       :", wrapped)
fmt.Println("Unwrap(wrapped):", errors.Unwrap(wrapped)) // lột 1 lớp → "tầng 2: user not found"
```

Output thật:

```text
wrapped       : tầng 1: tầng 2: user not found
Unwrap(wrapped): tầng 2: user not found
```

```text
  errors.Unwrap trả về         Khi nào nil?
  ─────────────────────────    ─────────────────────────────
  error cha TRỰC TIẾP (1 lớp)  error không wrap (%v/%s hoặc errors.New thuần)
```

> **Unwrap chỉ lột 1 lớp.** Muốn đi hết chuỗi → tự vòng lặp `for e := err; e != nil; e = errors.Unwrap(e)`.
> Nhưng thường ta dùng `Is`/`As` để duyệt tự động.

---

## 3. `errors.Is` — tìm theo **giá trị** (target)

`errors.Is(err, target)` duyệt **toàn bộ chuỗi**, trả `true` nếu gặp `target`:

```go
err := handlerGetUser(42) // đã wrap 2 tầng
fmt.Println("err chain          :", err)
fmt.Println("Is(ErrUserNotFound):", errors.Is(err, ErrUserNotFound))

// So sánh SAI cách: == trực tiếp thất bại vì đã wrap
fmt.Println("err == ErrUserNotFound:", err == ErrUserNotFound)
```

Output thật:

```text
err chain          : handler: service.GetUser(42): user not found
Is(ErrUserNotFound): true
err == ErrUserNotFound: false
```

```text
  err == target           errors.Is(err, target)
  ────────────────────    ─────────────────────────────
  Chỉ so SÁT lớp ngoài.   DUYỆT HẾT chuỗi (gọi Unwrap lặp).
  Wrap 1 lần là fail.      Đúng cả khi target nằm SÂU trong chuỗi.
```

### `errors.Is` so sánh như thế nào?

```text
errors.Is(err, target):
  1. Nếu err == target (== hoặc nil)          → true
  2. Nếu err có method Is(error) bool          → gọi err.Is(target)
  3. Nếu err có Unwrap() error                 → đệ quy Is(cha, target)
  4. Nếu err có Unwrap() []error (Go 1.20+)    → Is từng "cha" con
  5. Hết chuỗi mà không khớp                   → false
```

> Nhờ bước 2, bạn có thể **tự định nghĩa** cách so sánh bằng method `Is`
> (vd: so sánh theo mã lỗi thay vì identity) — xem mục 6.

---

## 4. `errors.As` — tìm theo **kiểu** (type)

`errors.As(err, &target)` duyệt chuỗi, tìm error có **kiểu concrete khớp** với `*target`,
rồi **gán** nó ra biến `target` — để lấy **dữ liệu structured**:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation %s: %s", e.Field, e.Message)
}

verr := fmt.Errorf("validate: %w", &ValidationError{Field: "age", Message: "phải >= 0"})

var vErr *ValidationError                    // biến nhận, kiểu TRÙNG với error gốc
if errors.As(verr, &vErr) {                  // &vErr phải là **con trỏ tới interface/kiểu**
    fmt.Printf("As: field=%s, message=%s\n", vErr.Field, vErr.Message)
}
```

Output thật:

```text
As: field=age, message=phải >= 0
```

```text
  errors.Is                          errors.As
  ────────────────────────────────   ──────────────────────────────────
  Tìm theo GIÁ TRỊ (target).         Tìm theo KIỂU (type).
  Quay lại bool (đúng/sai).          Gán error ra biến typed để DÙNG dữ liệu.
  Phù hợp: sentinel.                 Phù hợp: custom error struct.
```

### ⚠️ Lưu ý tham số `As`

* Tham số thứ 2 **phải là con trỏ** tới biến nhận (`&vErr`), và biến đó phải có kiểu
  **là interface hoặc concrete type** mà error trong chuỗi implement.
* Nếu chuỗi **không** chứa kiểu nào khớp → `As` trả `false`, biến không đổi.

---

## 5. Go 1.20+: nhiều `%w` & `errors.Join`

### Nhiều `%w` — một wrap có nhiều "cha"

```go
var ErrRead = errors.New("lỗi đọc")
var ErrParse = errors.New("lỗi parse")

err := fmt.Errorf("đọc file: %w (xử lý: %w)", ErrRead, ErrParse)
fmt.Println(errors.Is(err, ErrRead))   // true
fmt.Println(errors.Is(err, ErrParse))  // true
```

```text
Nhiều %w → error có Unwrap() []error → Is/As duyệt MỌI nhánh "cha".
```

### `errors.Join` — gộp nhiều error **độc lập**

```go
func validateUser(name, email string) error {
    var errs []error
    if name == "" {
        errs = append(errs, &ValidationError{Field: "name", Message: "không được để trống"})
    }
    if email == "" {
        errs = append(errs, &ValidationError{Field: "email", Message: "không được để trống"})
    }
    return errors.Join(errs...) // trả về 1 error chứa tất cả; nil nếu mọi input đều nil
}

if err := validateUser("", ""); err != nil {
    fmt.Println("Join gộp 2 lỗi:")
    if joined, ok := err.(interface{ Unwrap() []error }); ok {
        for _, e := range joined.Unwrap() {
            fmt.Println("  -", e)
        }
    }
}
```

Output thật:

```text
Join gộp 2 lỗi:
  - validation name: không được để trống
  - validation email: không được để trống
```

```text
  fmt.Errorf nhiều %w              errors.Join
  ─────────────────────────────    ────────────────────────────────
  Một CHUỖI wrap có nhiều cha.     Gộp NHIỀU lỗi độc lập (vd: validate).
  Message do bạn định dạng.        Message = nối các lỗi con (thêm \n).
```

---

## 6. Nâng cao: method `Is` / `Unwrap` tùy biến

### Method `Is` — tự định nghĩa cách so sánh

```go
type CodeError struct {
    Code int
}

func (e *CodeError) Error() string { return fmt.Sprintf("code %d", e.Code) }

// Is cho phép so sánh theo CODE thay vì identity (con trỏ).
func (e *CodeError) Is(target error) bool {
    t, ok := target.(*CodeError)
    return ok && t.Code == e.Code
}

err := &CodeError{Code: 404}
fmt.Println("Is theo code:", errors.Is(err, &CodeError{Code: 404})) // true — khác con trỏ nhưng cùng Code
```

```text
errors.Is gọi err.Is(target) → khớp theo Code, không cần cùng con trỏ.
```

### Method `Unwrap` — tham gia chuỗi

```go
type QueryError struct {
    SQL string
    Err error // lỗi gốc
}

func (e *QueryError) Error() string { return fmt.Sprintf("query %q: %v", e.SQL, e.Err) }

// Unwrap để errors.Is/As lọt tới e.Err (lỗi gốc).
func (e *QueryError) Unwrap() error { return e.Err }
```

> Implement `Unwrap() error` → custom error của bạn **trở thành một mắt xích** trong chuỗi,
> cho phép caller `errors.Is`/`errors.As` tới được lỗi gốc nằm bên trong.

---

## 7. Best Practices / Lưu ý quan trọng

* **Propagate bằng `%w`, không phải `%v`** — `%v` phá hủy cấu trúc, `Is`/`As` fail.
* **Wrap ngữ cảnh ở mỗi tầng** (`fmt.Errorf("tầng: %w", err)`) nhưng đừng wrap thừa.
* **Dùng `errors.Is`** để kiểm tra sentinel; **`errors.As`** để lấy dữ liệu từ custom error.
* **Không so sánh chuỗi** (`strings.Contains(err.Error(), ...)`) — dễ vỡ, chậm, thiếu chắc chắn.
* `errors.As` cần tham số **con trỏ** `&target`, biến nhận kiểu **trùng với error gốc**.
* Custom error muốn tham gia chain → implement **`Unwrap() error`**; muốn so sánh
  tùy biến → implement **`Is(error) bool`**.
* Go 1.20+: dùng **nhiều `%w`** hoặc **`errors.Join`** khi một thao tác vướng nhiều lỗi.
* `errors.Join` trả `nil` nếu mọi input đều nil — an toàn để accumulate lỗi rồi trả một lần.

---

## 8. Summary

* **`%w` wrap** lỗi gốc → tạo error chain; mỗi tầng thêm ngữ cảnh, giữ nguyên nguồn.
* **`errors.Unwrap`** lột **1 lớp**; **`errors.Is`** duyệt **hết chuỗi** tìm theo **giá trị**;
  **`errors.As`** duyệt tìm theo **kiểu** và **gán dữ liệu** ra biến typed.
* `errors.Is` so sánh bằng `==` hoặc method `Is` tùy biến; `errors.As` cần `&target`.
* Go 1.20+: **nhiều `%w`** (một wrap nhiều cha) và **`errors.Join`** (gộp lỗi độc lập).
* Custom error tham gia chain nhờ **`Unwrap() error`**; so sánh tùy biến nhờ **`Is(error) bool`**.
* Luôn propagate bằng `%w`; tuyệt đối không so sánh lỗi qua chuỗi `err.Error()`.
