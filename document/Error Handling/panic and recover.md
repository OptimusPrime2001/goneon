# `panic` và `recover` trong Go

`panic` và `recover` là cơ chế **dừng khẩn cấp** của Go — dùng cho tình huống
**bất thường không lường trước được**, **không** phải cách xử lý lỗi thay cho `error`.

Khớp node **Error Handling → Panic & Recover** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/panic_recover.go`

> ⚠️ **Nhắc lại triết lý Go:** lỗi **bình thường** → trả `error`.
> `panic` chỉ dành cho việc **mất invariant** khiến chương trình không thể tiếp tục an toàn.

---

## 1. Khái niệm & Bản chất

```text
  error (bình thường)                        panic (bất thường)
  ─────────────────────────────────          ──────────────────────────────────
  Là GIÁ TRỊ trả về.                         DỪNG chương trình khẩn cấp.
  Caller BẮT BUỘC kiểm tra.                  Unwind stack, chạy các defer đang đăng ký.
  Control flow tường minh.                   Không phải control flow để lạm dụng.
  Xử lý được, tiếp tục.                      Kết thúc goroutine/process nếu không recover.
```

**Khi `panic` xảy ra, Go:**
1. **Dừng** thực thi hàm hiện tại ngay lập tức.
2. **Chạy** các hàm `defer` **đã đăng ký** của các frame đang unwind (theo thứ tự **LIFO**).
3. Nếu **không có `recover()`** nào bắt → **crash** chương trình, in stack trace.

---

## 2. `panic` — dừng khẩn cấp, defer vẫn chạy

```go
func PanicDemo() {
    fmt.Println("Start")
    defer fmt.Println("Cleanup 1")
    defer fmt.Println("Cleanup 2") // chạy TRƯỚC Cleanup 1 (LIFO)

    panic("something went wrong")
    // Dòng dưới KHÔNG bao giờ được in ra.
    fmt.Println("End")
}
```

Output thật:

```text
Start
Cleanup 2
Cleanup 1
panic: something went wrong
... (stack trace, chương trình crash)
```

```text
  Thứ tự chạy khi panic:
  ─────────────────────────────────────────────
  Start
  → panic xảy ra, dừng tại đây
  → Cleanup 2   (defer đăng ký SAU, chạy TRƯỚC — LIFO)
  → Cleanup 1   (defer đăng ký TRƯỚC, chạy SAU)
  → CRASH (vì không có recover)
```

> **Điểm mấu chốt:** `defer` **vẫn chạy** dù `panic` xảy ra — đây là nền tảng
> để `recover()` hoạt động (nó nằm trong một hàm `defer`).

---

## 3. `recover` — bắt panic trong `defer`

`recover()` **chỉ có tác dụng khi gọi bên trong hàm `defer`**. Nó "bắt" panic,
chuyển trạng thái goroutine về bình thường, chương trình **sống tiếp**:

```go
func RecoverDemo() {
    defer func() {
        // recover() trả giá trị panic nếu có, ngược lại nil.
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    fmt.Println("Before panic")
    panic("something went wrong")

    // KHÔNG được thực thi (panic xảy ra trước).
    fmt.Println("After panic")
}
```

Output thật:

```text
Before panic
Recovered: something went wrong
```

```text
  recover() gọi NGOÀI defer → trả nil, KHÔNG bắt được panic.
  recover() trả nil khi KHÔNG có panic đang diễn ra.
  → Luôn kiểm tra: if r := recover(); r != nil { ... }
```

---

## 4. Panic từ runtime & chuyển panic → error

Một số thao tác **runtime panic** khi dữ liệu sai (không phải do bạn `panic` tường minh):

```go
// DivideSafe bọc hàm prone-to-panic rồi recover, trả về error bình thường.
// Đây là MÔ HÌNH chuẩn: panic → recover → chuyển thành error.
func DivideSafe(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("DivideSafe: %v", r) // chuyển panic thành error
        }
    }()
    return a / b, nil // panic runtime "integer divide by zero" nếu b == 0
}
```

```go
if r, err := DivideSafe(10, 0); err != nil {
    fmt.Println("DivideSafe lỗi:", err)
} else {
    fmt.Println("DivideSafe =", r)
}
```

Output thật:

```text
DivideSafe lỗi: DivideSafe: runtime error: integer divide by zero
```

```text
  Runtime panic thường gặp (không do bạn chủ động):
  ─────────────────────────────────────────────────────────────
  • integer divide by zero          a / 0
  • index out of range              arr[10] với len(arr)=5
  • assignment to entry in nil map  m[k] = v với m == nil
  • nil pointer dereference         *p với p == nil
```

> **Mô hình chuẩn:** nếu bạn wrap một hàm **dễ panic** (vd: chia cho 0, parse),
> dùng `defer + recover` để **chuyển panic thành `error`** — caller nhận error bình thường.

---

## 5. Khi nào panic hợp lệ? (và khi nào thì KHÔNG)

```text
  NÊN panic (mất invariant nghiêm trọng)      KHÔNG NÊN panic (dùng error)
  ────────────────────────────────────────    ─────────────────────────────────
  • Init thất bại, app không thể chạy.        • Lỗi đầu vào do user gửi.
  • Vi phạm invariant nội bộ (bug code).      • File không tồn tại.
  • Giao thức bị phá vỡ, không thể tiếp tục.  • Lỗi mạng, DB timeout.
  • Trong test (t.Fatalf thay vì panic).       • Bất kỳ tình huống "bình thường" nào.
```

```text
  Nguyên tắc chung:
  ─────────────────────────────────────────────────────────────
  • Library KHÔNG NÊN panic cho lỗi đầu vào → trả error.
  • panic gần giống "assert": dùng khi TIẾP TỤC SẼ NGUY HIỂM hơn là dừng.
  • Nếu bạn có thể xử lý và trả về được → trả error, đừng panic.
```

---

## 6. Best Practices / Lưu ý quan trọng

* **`panic`/`recover` KHÔNG phải cơ chế xử lý lỗi thay cho `error`.**
  Lỗi bình thường → trả `error`; panic chỉ cho mất invariant không thể phục hồi.
* **`recover()` phải nằm trong hàm `defer`** — gọi ở nơi khác vô tác dụng.
* Luôn kiểm tra `if r := recover(); r != nil` — recover trả `nil` khi không có panic.
* **Recover ở tầng boundary** (HTTP handler, worker goroutine, main) — không recover
  rải rác khắp nơi (che bug, làm code rối).
* **Ghi log stack trace khi recover** — đừng nuốt im lặng (xem [Stack Traces](<Stack Traces and Debugging.md>)).
* Library **trả `error`**, không panic cho lỗi đầu vào.
* Mô hình tốt: wrap hàm prone-to-panic bằng `defer + recover`, **chuyển panic → error**.

---

## 7. Summary

* `panic` **dừng khẩn cấp**, **unwind stack**, **chạy defer** (LIFO), crash nếu không recover.
* `recover()` **chỉ hiệu lực trong hàm `defer`**; trả `nil` khi không có panic.
* **Runtime panic** (chia 0, index out of range, nil map, nil pointer) — dữ liệu sai.
* Mô hình chuẩn: `defer + recover` để **chuyển panic thành `error`** cho caller.
* Panic hợp lệ cho **mất invariant nghiêm trọng**; lỗi bình thường → **trả `error`**.
* Recover ở **tầng boundary**, ghi log stack trace, đừng recover mù quáng khắp nơi.
