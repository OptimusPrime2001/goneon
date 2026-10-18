# Stack Traces & Debugging trong Go

Stack trace là **"dấu vân tay"** cho biết chương trình đang ở đâu khi xảy ra lỗi/panic.
Hiểu cách đọc và lấy stack trace giúp bạn debug nhanh, đặc biệt khi recover panic.

Khớp node **Error Handling → Stack Traces & Debugging** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/stack_trace.go`

---

## 1. Khái niệm & Bản chất

```text
  Stack trace = danh sách các FRAME (khung gọi) đang hoạt động, mỗi frame gồm:
  ────────────────────────────────────────────────────────────────────────────
  • Tên hàm   (main.riskyFunc)
  • File      (stack_trace.go:10)
  • Dòng code (+0xaa = offset byte, hiếm khi cần)
  • Goroutine (goroutine 1 [running])
```

Stack trace sinh ra khi:
* **Panic** xảy ra (runtime tự in ra).
* Bạn chủ động gọi `runtime/debug.Stack()` (thường khi recover).
* Lỗi fatal không phục hồi được.

> **Bản chất:** Go **không** nhồi stack trace vào `error` như một số ngôn ngữ khác.
> Nếu muốn giữ trace, bạn phải **tự gắn** nó (vd: khi recover panic).

---

## 2. Lấy stack trace: `runtime/debug.Stack()`

`debug.Stack()` trả về **chuỗi byte** là stack trace của **goroutine hiện tại** tại thời điểm gọi:

```go
func riskyFunc() {
    panic("đã xảy ra lỗi nghiêm trọng")
}

func middleLayer() { riskyFunc() } // frame trung gian

func StackTrace() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
            fmt.Println("---- Stack trace (debug.Stack) ----")
            fmt.Println(string(debug.Stack()))
        }
    }()
    middleLayer()
}
```

Output thật (rút gọn, path cụ thể tùy máy):

```text
Recovered: đã xảy ra lỗi nghiêm trọng
---- Stack trace (debug.Stack) ----
goroutine 1 [running]:
runtime/debug.Stack()
        .../src/runtime/debug/stack.go:26 +0x5e
main.StackTrace.func1()
        .../stack_trace.go:26 +0xaa
panic({0x...?, 0x...?})
        .../src/runtime/panic.go:860 +0x13a
main.riskyFunc(...)
        .../stack_trace.go:10
main.middleLayer(...)
        .../stack_trace.go:15
main.StackTrace()
        .../stack_trace.go:30 +0x45
main.main()
        .../main.go:25 +0x245
```

---

## 3. Đọc stack trace — đọc từ dưới lên

```text
  Cấu trúc 1 frame:
  ─────────────────────────────────────────────
  main.riskyFunc(...)
          .../stack_trace.go:10
          ▲         ▲            ▲
          │         │            └── dòng code trong file
          │         └── file chứa hàm
          └── tên hàm (package.Func)
```

**Cách đọc:** bắt đầu từ **dưới cùng** (gần `main` nhất) và đi **lên trên**
(gọi gần nhất, nơi panic xảy ra):

```text
  main.main()            ← điểm bắt đầu (gọi đầu tiên)
      ↓
  main.StackTrace()      ← tầng cao
      ↓
  main.middleLayer()     ← tầng giữa
      ↓
  main.riskyFunc()       ← tầng sâu nhất — NƠI PANIC XẢY RA
      ↓
  panic(...)             ← chính panic
```

```text
  Dấu (...) sau tên hàm có nghĩa: hàm đó đã INLINE (nội suy) vào hàm gọi.
  Gợi ý đọc trace:
  • Frame trên CÙNG (debug.Stack, panic) là runtime/framework — thường bỏ qua.
  • Frame ĐẦU TIÊN của code bạn (từ dưới lên) = nơi lỗi GẦN NHẤT.
  • Đọc ngược lên để thấy LUỒNG GỌI (tầng nào gọi tầng nào).
```

---

## 4. Tăng độ chi tiết traceback: `GOTRACEBACK`

Môi trường `GOTRACEBACK` (hoặc `debug.SetTraceback`) kiểm soát **mức chi tiết** của stack trace:

```text
  GOTRACEBACK         Kết quả
  ─────────────────   ──────────────────────────────────────────────
  "single" (mặc định) Chỉ trace goroutine gây panic, bỏ qua system goroutines.
  "all"               Trace TẤT CẢ goroutine đang chạy — hữu ích khi deadlock/race.
  "system"            Như "all" + thêm thông tin runtime.
  "crash"             Crash thật sự (tạo core dump) để phân tích sâu.
```

```go
// Trong code: đặt TRƯỚC khi panic (thường trong main/init)
debug.SetTraceback("all")
```

```bash
# Hoặc qua biến môi trường khi chạy:
GOTRACEBACK=all go run ./...
```

> Dùng `GOTRACEBACK=all` khi program **treo/deadlock** — nó in trace của
> **mọi goroutine** để bạn thấy goroutine nào đang chờ ở đâu.

---

## 5. Gắn stack trace vào error khi recover

Thay vì chỉ in ra, bạn có thể **giữ** trace trong error:

```go
func SafeHandler() (err error) {
    defer func() {
        if r := recover(); r != nil {
            // debug.Stack() trả []byte → chuyển sang string để nhúng vào error
            err = fmt.Errorf("panic: %v\n%s", r, string(debug.Stack()))
        }
    }()
    riskyFunc() // có thể panic
    return nil
}
```

> ⚠️ **Chỉ làm vậy khi cần** (vd: log phía server để truy vết).
  `debug.Stack()` **tốn kém** — không dùng trong code path thường xuyên.

---

## 6. Mẹo debug lỗi (không panic) nhanh

```text
  Mẹo                               Áp dụng
  ────────────────────────────────  ──────────────────────────────────────────
  In error với %+v                   fmt.Printf("%+v", err) — thấy thêm ngữ cảnh wrap
  Wrap + in chuỗi lỗi               errors.Is/As để lật chuỗi, in từng lớp
  Delve (dlv)                        Debugger chính thức: breakpoint, print, stack
  go vet / staticcheck               Bắt lỗi tĩnh (unused err, printf verb sai)
  GOTRACEBACK=all                    In trace MỌI goroutine khi treo/deadlock
```

```go
// %+v trên error thường cho thêm chi tiết (tùy type error)
err := fmt.Errorf("handler: %w", ErrUserNotFound)
fmt.Printf("v : %v\n", err)   // handler: user not found
fmt.Printf("+v: %+v\n", err)  // handler: user not found  (type hỗ trợ %+v thì có thêm field)
```

---

## 7. Best Practices / Lưu ý quan trọng

* Stack trace Go đọc từ **dưới lên** — frame code bạn đầu tiên (sau runtime) là nơi lỗi gần nhất.
* **Không nhồi `debug.Stack()` vào mọi error** — chỉ khi recover panic / log server-side.
* `debug.Stack()` lấy trace của **goroutine hiện tại** — gọi đúng chỗ cần (trong `defer` recover).
* Dùng **`GOTRACEBACK=all`** khi cần thấy **mọi goroutine** (deadlock, race, treo).
* Error Go **không mặc định kèm stack trace** — nếu muốn giữ trace, tự gắn lúc recover.
* Kết hợp `defer + recover + debug.Stack()` để **biến panic thành error có trace**.
* Dùng **Delve (`dlv`)** để debug chủ động (breakpoint, inspect) thay vì chỉ in log.

---

## 8. Summary

* **Stack trace** = danh sách frame (hàm, file, dòng) đang hoạt động khi xảy ra lỗi/panic.
* **`debug.Stack()`** trả về trace của **goroutine hiện tại** — dùng khi recover panic.
* Đọc trace từ **dưới lên** (từ `main` tới nơi panic); dấu `(...)` = hàm bị inline.
* **`GOTRACEBACK`** (`single`/`all`/`system`/`crash`) điều chỉnh độ chi tiết; `all` cho mọi goroutine.
* Error Go **không tự kèm trace** — muốn giữ thì **tự gắn** lúc recover.
* Mẹo: `%+v` cho thêm chi tiết, `Delve` để debug chủ động, tránh lạm dụng `debug.Stack()`.
