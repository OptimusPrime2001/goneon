# Channels & Select trong Go

Channel là cách các goroutine **truyền dữ liệu cho nhau** — thay vì sửa chung
một biến (gây data race).

Code mẫu: `fundamental/5.concurrency/channel.go`
Đọc trước: [Goroutines](<Goroutines.md>)

---

## 1. Channel là gì?

```go
ch := make(chan int)     // unbuffered — không có hàng đợi
ch2 := make(chan int, 3)  // buffered  — chứa được tối đa 3 phần tử
```

```text
  "Don't communicate by sharing memory;
   share memory by communicating."
   → Thay vì 2 goroutine cùng sửa 1 biến, 1 cái GỬI, 1 cái NHẬN.
   → Dữ liệu chỉ có MỘT chủ tại mỗi thời điểm → không cần Mutex.
```

---

## 2. Unbuffered vs Buffered

```text
  UNBUFFERED (cap = 0)                BUFFERED (cap = n)
  ─────────────────────────           ─────────────────────────
  Người gửi ĐỨNG CHỜ có người nhận   Gửi vào hàng đợi rồi đi tiếp
  → đồng bộ, tay bắt tay              Không cần người nhận ngay
                                      Hàng ĐẦY → người gửi block
```

```go
func BufferedDemo() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    // ch <- 3     // uncomment → block: hàng đợi đầy, không ai đọc

    fmt.Println(len(ch), cap(ch))   // 2 2
    fmt.Println(<-ch)               // 1
    fmt.Println(len(ch), cap(ch))   // 1 2
}
```

```text
  len = số phần tử đang chờ, cap = sức chứa.
  Unbuffered là "đồng bộ THẬT": không Sleep, không đoán —
  chính channel tạo điểm gặp nhau.
```

---

## 3. Directional — ràng buộc ai gửi, ai nhận

```go
func producer(out chan<- int) {   // chỉ được GỬI
    for i := 1; i <= 3; i++ {
        out <- i * 10
    }
    close(out)                    // người gửi close
}

func consumer(in <-chan int) {    // chỉ được NHẬN
    for v := range in {           // range tự dừng khi close
        fmt.Println("nhận:", v)   // 10, 20, 30
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)
    consumer(ch)
}
```

```text
  • chan<- T  = send-only   ,  <-chan T = receive-only.
  • SAI chiều → LỖI BIÊN DỊCH (không phải runtime).
  • Chữ ký hàm nói rõ luôn ai là producer, ai là consumer.
```

---

## 4. `close()` & `range`

**Người GỬI close** — không bao giờ để người nhận close.

```go
ch := make(chan string, 3)
ch <- "a"; ch <- "b"; ch <- "c"
close(ch)

for v := range ch {   // lặp tới khi channel đóng rồi thoát
    fmt.Println(v)    // a, b, c
}
```

```text
  • close TRÊN CHANNEL ĐÃ ĐÓNG   → PANIC.
  • close TRÊN nil               → PANIC.
  • KHÔNG close mà range          → block mãi → goroutine leak.
  • Đóng rồi vẫn đọc được phần tử còn trong hàng đợi, hết thì trả zero value.
```

### Phân biệt "đóng" với "giá trị 0"

```go
ch <- 7
close(ch)

v, ok := <-ch    // v=7,  ok=true   → còn dữ liệu
v, ok = <-ch     // v=0,  ok=false  → channel đã đóng
```

> Cần vì `0` cũng là giá trị hợp lệ của `chan int` — nhìn `v` không phân biệt được.

---

## 5. Select statement

Chọn chạy **case nào sẵn sàng trước**. Nhiều case sẵn sàng → chọn **NGẪU NHIÊN**.

```go
func SelectDemo() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() { time.Sleep(50 * time.Millisecond); c1 <- "từ c1" }()
    go func() { time.Sleep(100 * time.Millisecond); c2 <- "từ c2" }()

    for i := 0; i < 2; i++ {
        select {
        case msg := <-c1:
            fmt.Println("nhận:", msg)
        case msg := <-c2:
            fmt.Println("nhận:", msg)
        }
    }
}
```

Output:

```text
nhận: từ c1
nhận: từ c2
```

```text
  • KHÔNG case nào sẵn sàng → select BLOCK tới khi có.
  • case GỬI (ch <- v) cũng được như case nhận.
  • Cú pháp như switch nhưng KHÔNG có expression: select { ... }
```

### 5.1 `default` — không bao giờ block

```go
select {
case v := <-ch:
    fmt.Println("nhận:", v)
default:
    fmt.Println("không có gì sẵn → chạy ngay, không chờ")
}
```

> Pattern "thử lấy, không có thì thôi" — đừng dựa vào `len(ch)`.

### 5.2 Timeout

```go
select {
case v := <-ch:
    fmt.Println("nhận:", v)
case <-time.After(100 * time.Millisecond):
    fmt.Println("timeout: không có dữ liệu sau 100ms")
}
```

> Trong vòng lặp dày đặc → dùng `time.NewTimer` + `Stop()` thay `time.After` (tránh rò rỉ timer).

### 5.3 Channel `nil` trong select

```go
var ch chan int     // nil

select {
case v := <-ch:
    fmt.Println("KHÔNG bao giờ chạy", v)
case <-time.After(50 * time.Millisecond):
    fmt.Println("case nil bị bỏ qua → timeout thắng")
}
```

```text
  • nil channel NGOÀI select → block MÃI (goroutine leak).
  • nil channel TRONG select → case đó bị bỏ qua (rất hữu ích để "tắt" nhánh).
```

---

## 6. Lỗi thường gặp

| Hành vi | Kết quả |
|---------|---------|
| Gửi trên channel **đã close** | **Panic**: `send on closed channel` |
| `close()` 2 lần / close `nil` | **Panic** |
| `<-ch` với `ch == nil` (ngoài select) | **Block mãi** |
| `range` channel không close | **Block mãi** → leak |
| `select {}` không có case | **Block mãi** (biên dịch vẫn OK) |
| Đọc channel đã đóng, hết phần tử | Trả zero value, `ok=false` — **không panic** |

```go
ch := make(chan int)
ch <- 1          // không ai đọc → block
// → deadlock: all goroutines are asleep - deadlock!
```

```text
  Go chỉ phát hiện được deadlock khi TẤT CẢ goroutine đều block.
  Deadlock nửa chương → Go không báo, chỉ tốn RAM + goroutine leak.
```

---

## 7. Best Practices

* **Người gửi close**, close xong đừng gửi nữa (sai = panic).
* **Buffered nhỏ có lý do** (1, số worker) — đừng đệm lớn "cho chắc" (che vấn đề).
* Dùng **directional channel** trong chữ ký hàm để ghi rõ contract.
* **`nil` channel = block mãi** — kiểm tra channel đã được `make`/gán trước khi dùng.
* Channel **không thay thế Mutex** — chỉ cần bảo vệ dữ liệu thì dùng [Synchronization](<Synchronization.md>).
* Chạy `go test -race ./...` khi có nhiều goroutine + channel.

---

## 8. Summary

* Unbuffered = **đồng bộ tay bắt tay**; buffered = hàng đợi, đầy thì block.
* **Người gửi `close()`**; `range` tự thoát; `v, ok := <-ch` phân biệt đóng vs zero value.
* **`select`** chọn case sẵn sàng, `default` = non-block, `time.After` = timeout,
  `case nil` = bỏ qua nhánh.
* Sai close → **panic**; không close → **leak**; `select {}` → **block mãi**.

---

## Bài tiếp theo

| Bài | Nội dung |
|-----|-----------|
| ~~Select~~ | ❌ Đã gộp vào Channels.md (mục 5) |
| [Worker Pool](<Worker Pool.md>) ⬜ | `sync.WaitGroup`, pool cố định, job động, chống deadlock |
| [Synchronization](<Synchronization.md>) ⬜ | Mutex, RWMutex, Atomic |
| [Context](<Context.md>) ⬜ | Cancel, timeout |
| [Concurrency Patterns](<Concurrency Patterns.md>) ⬜ | Fan-Out/In, Pipeline, Producer-Consumer |
