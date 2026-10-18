# Tạo Error: `errors.New` & `fmt.Errorf`

Go cung cấp 2 cách tạo error phổ biến: `errors.New` cho lỗi **nguyên sơ, tĩnh**,
và `fmt.Errorf` cho lỗi **có ngữ cảnh, có thể format** (đặc biệt là `%w` để wrap).

Khớp node **Error Handling → errors.New / fmt.Errorf** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/error_basic.go`

---

## 1. Khái niệm & Bản chất

```text
  errors.New(msg)                  fmt.Errorf(format, args...)
  ───────────────────────────      ─────────────────────────────────────
  Trả về error CỐ ĐỊNH.           Format chuỗi rồi bọc thành error.
  Message nguyên sơ, không đổi.    Hỗ trợ %v / %s (chuỗi) và %w (wrap).
  Phù hợp: sentinel, lỗi tĩnh.    Phù hợp: thêm ngữ cảnh, chèn biến.
```

Cả hai đều trả về interface `error` (xem [The error Interface](<The error Interface.md>)).

---

## 2. `errors.New` — error nguyên sơ

`errors.New` tạo một `error` với message **cố định** (static):

```go
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("divide: chia cho 0")
    }
    return a / b, nil
}
```

```text
Bản chất: errors.New trả về *errors.errorString — struct PRIVATE chỉ có 1 field
message string, và method Error() trả về chính message đó.
```

> **Khi nào dùng `errors.New`?**
> * Tạo **sentinel error** cho package (bài [Sentinel Errors](<Sentinel Errors.md>)).
> * Lỗi **ngắn gọn, không cần chèn biến** (vd: `"EOF"`, `"user not found"`).

### ⚠️ Không gọi `errors.New` ở "điều kiện bình thường"

`errors.New` **phân bổ bộ nhớ** (alloc). Tạo nó lặp lại mỗi lần gọi hàm là lãng phí:

```go
// ❌ SAI — tạo error mới MỖI lần, kể cả khi không lỗi
func DivideBad(a, b int) (int, error) {
    err := errors.New("divide: chia cho 0") // alloc vô ích khi b != 0
    if b == 0 {
        return 0, err
    }
    return a / b, nil
}

// ĐÚNG — chỉ tạo khi THỰC SỰ gặp lỗi
func DivideGood(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("divide: chia cho 0") // alloc đúng lúc cần
    }
    return a / b, nil
}
```

> **Quy tắc:** `errors.New`/`fmt.Errorf` chỉ gọi trong **nhánh lỗi**,
> không gọi trước `if` rồi để đó "phòng hờ".

---

## 3. `fmt.Errorf` — format + ngữ cảnh

`fmt.Errorf` hoạt động như `fmt.Sprintf` nhưng **trả về `error`** thay vì `string`:

```go
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // Chèn biến (path) vào message để biết lỗi xảy ra ở đâu
        return nil, fmt.Errorf("LoadConfig(%s): %w", path, err)
    }
    _ = data
    return &Config{}, nil
}
```

```text
So sánh cú pháp verb:
  %s / %v   →  "chuỗi hóa" error thành message, MẤT cấu trúc (không unwrap được)
  %w        →  "wrap" error gốc, GIỮ NGUYÊN cấu trúc (unwrap/Is/As được)
```

---

## 4. `%w` vs `%v`/`%s` — khác biệt cốt lõi

Đây là điểm **quan trọng nhất** khi dùng `fmt.Errorf`.

```go
origin := errors.New("nguồn lỗi")

// %v / %s: chỉ copy CHUỖI message, mất đối tượng error gốc
byValue := fmt.Errorf("tầng A: %v", origin)

// %w: bọc error gốc, GIỮ NGUYÊN object → unwrap được
byWrap := fmt.Errorf("tầng A: %w", origin)

fmt.Println("byValue:", byValue)
fmt.Println("byWrap :", byWrap)

// So sánh cấu trúc
fmt.Println("Unwrap(byValue) == nil ?", errors.Unwrap(byValue) == nil) // true — mất gốc
fmt.Println("Unwrap(byWrap)  == origin ?", errors.Unwrap(byWrap) == origin) // true — còn gốc
```

Output thật:

```text
byValue: tầng A: nguồn lỗi
byWrap : tầng A: nguồn lỗi
Unwrap(byValue) == nil ? true
Unwrap(byWrap)  == origin ? true
```

```text
  Verb   Message hiển thị   Giữ error gốc?   errors.Is/As hoạt động?
  ─────  ─────────────────  ───────────────  ──────────────────────
  %s/%v  Có (chuỗi)         ❌ KHÔNG          ❌ KHÔNG
  %w     Có (chuỗi)         ✓ CÓ             ✓ CÓ
```

> **Hậu quả thực tế:** dùng `%v` khi propagate → tầng trên gọi `errors.Is(err, ErrNotFound)`
> sẽ **fail** (trả false) vì error gốc đã bị "phẳng hóa" thành chuỗi.
> → **Luôn dùng `%w`** khi bọc lỗi để propagate.

### Nhiều `%w` trong 1 lệnh (Go 1.20+)

Từ Go 1.20, `fmt.Errorf` cho phép **nhiều `%w`** — error bọc **nhiều "cha"** cùng lúc
(`Unwrap()` trả về slice). Phù hợp khi một thao tác vướng nhiều lỗi:

```go
// Dùng sentinel để so sánh được (errors.New mỗi lần tạo object MỚI → == false)
var ErrRead = errors.New("lỗi đọc")
var ErrParse = errors.New("lỗi parse")

err := fmt.Errorf("đọc file: %w (xử lý: %w)", ErrRead, ErrParse)
fmt.Println("message:", err)

// errors.Is tìm thấy BẤT KỲ "cha" nào trong tree
fmt.Println("Is(ErrRead):", errors.Is(err, ErrRead))
fmt.Println("Is(ErrParse):", errors.Is(err, ErrParse))
```

Output thật:

```text
message: đọc file: lỗi đọc (xử lý: lỗi parse)
Is(ErrRead): true
Is(ErrParse): true
```

> **So với `errors.Join`:** `Join` gộp các error **độc lập** (vd: lỗi validate nhiều field),
> còn nhiều `%w` dùng khi **một chuỗi wrap** cần giữ nhiều "cha". Chi tiết: [Error Wrapping](<Error Wrapping.md>).
> ⚠️ Muốn `errors.Is` đúng, phải so sánh với **cùng một biến sentinel**, không gọi
> `errors.New` lại (mỗi lần là object mới → `==` luôn false).

---

## 5. Best Practices / Lưu ý quan trọng

* **`errors.New`**: message tĩnh, ngắn gọn; dùng cho sentinel hoặc lỗi cố định.
* **`fmt.Errorf`**: chèn biến + ngữ cảnh; dùng `%w` (không phải `%v`) khi propagate.
* **Chỉ tạo error trong nhánh lỗi** — đừng `errors.New` trước `if` cho "phòng hờ" (lãng phí alloc).
* **`%w` giữ nguyên error gốc** để `errors.Is`/`errors.As`/`errors.Unwrap` hoạt động;
  `%v`/`%s` "phẳng hóa" thành chuỗi, mất cấu trúc.
* **Message lower-case, không dấu chấm cuối**, thêm ngữ cảnh "đang làm gì" ở mỗi tầng.
* Go 1.20+ cho phép **nhiều `%w`** — dùng khi một thao tác gặp nhiều lỗi con (hoặc `errors.Join`).

---

## 6. Summary

* `errors.New(msg)` → error **tĩnh, nguyên sơ**, alloc trong nhánh lỗi.
* `fmt.Errorf(format, ...)` → error **format được**, chèn biến + ngữ cảnh.
* **Khác biệt then chốt:** `%w` **giữ error gốc** (unwrap/Is/As hoạt động),
  `%v`/`%s` chỉ copy chuỗi (**mất cấu trúc**).
* Go 1.20+: cho phép **nhiều `%w`** trong 1 `fmt.Errorf` (error tree).
* Quy tắc: propagate bằng `%w`, viết message lower-case, wrap ngữ cảnh mỗi tầng.
