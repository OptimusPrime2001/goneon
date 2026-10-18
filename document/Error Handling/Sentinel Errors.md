# Sentinel Errors trong Go

Sentinel error là một error **danh nghĩa** được khai báo sẵn dưới dạng **biến package-level**,
để caller có thể **so sánh và phân nhánh** theo nó.

Khớp node **Error Handling → Sentinel Errors** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/sentinel.go`

---

## 1. Khái niệm & Bản chất

**Sentinel** = "biến kiểm chứng" — một giá trị cố định mà các hàm **trả về** khi gặp
một tình huống cụ thể, và caller **kiểm tra** bằng `errors.Is`:

```go
// Khai báo package-level, export (chữ hoa đầu) để caller khác dùng được.
var ErrNotFound = errors.New("item not found")

// Hàm trả về sentinel khi không tìm thấy.
func FindItem(id int) (string, error) {
    if id != 1 {
        return "", ErrNotFound // trả về CHÍNH biến sentinel, không tạo mới
    }
    return "apple", nil
}
```

```text
  errors.New(...) mỗi lần gọi            Sentinel (biến package-level)
  ──────────────────────────────────     ──────────────────────────────────
  Tạo object MỚI mỗi lần.                Tạo MỘT LẦN, dùng lại.
  == / Is không so sánh được (khác       So sánh được bằng errors.Is (cùng
  con trỏ).                              biến → identity).
  Phù hợp: lỗi chi tiết, có dữ liệu.    Phù hợp: "cờ" danh nghĩa để phân nhánh.
```

> **Bản chất:** sentinel là một error **cố định, duy nhất** trong package.
> Các tầng wrap nó bằng `%w` nhưng **vẫn trỏ về cùng biến** → caller `errors.Is` được.

---

## 2. So sánh `==` vs `errors.Is`

```go
wrapped := fmt.Errorf("cache lookup: %w", ErrNotFound)
fmt.Println("Is(wrapped, ErrNotFound):", errors.Is(wrapped, ErrNotFound))
fmt.Println("== trực tiếp           :", wrapped == ErrNotFound)
```

Output thật:

```text
Is(wrapped, ErrNotFound): true
== trực tiếp           : false
```

```text
  wrapped == ErrNotFound        errors.Is(wrapped, ErrNotFound)
  ──────────────────────────    ────────────────────────────────────
  So sánh SÁT lớp ngoài.        Duyệt hết chuỗi (Unwrap lặp).
  Đã wrap 1 lần → FALSE.        Wrap bao nhiêu lần vẫn TRUE.
```

> **Luôn dùng `errors.Is`** thay vì `==` — vì error rất dễ bị wrap ở tầng giữa
> mà caller không kiểm soát được.

---

## 3. Phân nhánh theo sentinel

Sentinel biến việc kiểm tra lỗi thành **control flow bình thường**:

```go
if _, err := FindItem(99); err != nil {
    switch {
    case errors.Is(err, ErrNotFound):
        fmt.Println("Xử lý riêng: tạo mặc định hoặc 404")
    case errors.Is(err, ErrForbidden):
        fmt.Println("Xử lý riêng: 403")
    default:
        fmt.Println("Lỗi không xác định:", err)
    }
}
```

Output thật:

```text
Xử lý riêng: tạo mặc định hoặc 404
```

> Pattern `switch { case errors.Is(...): ... }` rất phổ biến — mỗi sentinel ứng với
> **một nhánh xử lý riêng**.

---

## 4. Sentinel có sẵn của stdlib

Go chuẩn đã định nghĩa sẵn nhiều sentinel quen thuộc:

```go
fmt.Println("io.EOF       :", io.EOF)
fmt.Println("os.ErrNotExist:", os.ErrNotExist)

if _, err := os.Open("khong_ton_tai.txt"); errors.Is(err, os.ErrNotExist) {
    fmt.Println("Is(err, os.ErrNotExist): true")
}
```

Output thật:

```text
io.EOF       : EOF
os.ErrNotExist: file does not exist
Is(err, os.ErrNotExist): true
```

```text
  Sentinel             Ý nghĩa
  ──────────────────   ────────────────────────────────────────────
  io.EOF               Đọc hết dữ liệu (không phải lỗi nghiêm trọng)
  io.ErrUnexpectedEOF  Hết dữ liệu đột ngột (có thật là lỗi)
  os.ErrNotExist       File/đường dẫn không tồn tại
  fs.ErrNotExist       (Go 1.20+) — bản dùng chung cho filesystem
  sql.ErrNoRows        Query trả 0 dòng (khi Scan 1 dòng)
  context.Canceled     Context bị hủy
  context.DeadlineExceeded  Context hết thời gian
```

> **Lưu ý:** `io.EOF` là sentinel **đặc biệt** — nó báo "hết dữ liệu", thường được
> coi là **kết thúc bình thường** chứ không phải lỗi cần xử lý nghiêm trọng.

---

## 5. Sentinel vs Custom Error — khi nào dùng cái nào?

```text
  Dùng SENTINEL khi                      Dùng CUSTOM ERROR (struct) khi
  ──────────────────────────────────     ─────────────────────────────────────
  • Caller cần PHÂN NHÁNH theo lỗi.     • Lỗi mang theo DỮ LIỆU (Code, Field...).
  • Lỗi đơn giản, "cờ" danh nghĩa.      • Caller cần trích xuất chi tiết (As).
  • Vd: ErrNotFound, ErrForbidden.      • Vd: *ValidationError, *QueryError.
```

Cả hai đều **wrap bằng `%w`** và kiểm tra bằng `errors.Is`/`errors.As`.

---

## 6. Best Practices / Lưu ý quan trọng

* Sentinel nên là **biến package-level export** (`ErrXxx`) để caller import và so sánh.
* Khai báo tập trung trong `var (...)`, dùng `errors.New` (message ngắn, ổn định).
* **Kiểm tra bằng `errors.Is`**, không dùng `==` (vì lỗi có thể đã wrap).
* **Không match sentinel qua chuỗi message** (`err.Error() == "..."`) — dễ vỡ khi đổi message.
* Chỉ dùng sentinel cho trường hợp **caller thật sự cần phân nhánh**;
  nếu cần dữ liệu chi tiết → custom error struct (dùng `errors.As`).
* Sentinel là **"hợp đồng" public** — đổi message không sao, nhưng **đổi tên biến**
  là breaking change cho caller.

---

## 7. Summary

* **Sentinel** = error danh nghĩa là **biến package-level export**, dùng `errors.New`.
* Hàm trả về **chính biến sentinel** (không tạo mới) khi gặp tình huống ứng với nó.
* **Kiểm tra bằng `errors.Is`** (duyệt hết chuỗi), **không** dùng `==` (chỉ so lớp ngoài).
* Go stdlib có sẵn: `io.EOF`, `os.ErrNotExist`, `sql.ErrNoRows`, `context.Canceled`...
* Sentinel cho **phân nhánh**; custom error struct cho **mang dữ liệu** (dùng `errors.As`).
* Sentinel là "hợp đồng" public — **không** match qua chuỗi message, chỉ qua `errors.Is`.
