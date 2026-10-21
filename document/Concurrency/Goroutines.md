# Goroutines trong Go

Goroutine là đơn vị thực thi **nhẹ** của Go — một hàm đang chạy **độc lập**,
được Go scheduler điều phối trên các OS thread.

Khớp node **Concurrency → Goroutines** của roadmap.sh (§11).
Code mẫu: `fundamental/5.concurrency/goroutines.go`

> ⚠️ **Bài này chỉ đi phần nền tảng:** tạo goroutine, lifecycle, đồng bộ sơ bộ bằng
> channel, data race, goroutine leak, các pitfall thường gặp.
> `WaitGroup`, `Mutex`, `Atomic`, `Context`, Worker Pool → đọc ở các bài sau.

---

## 1. Goroutine là gì?

```text
  OS Thread (heavy)                          Goroutine (light)
  ──────────────────────────────────         ─────────────────────────────────
  ~1MB stack (cấp phát bởi OS).              ~2-8KB stack (cấp phát bởi runtime,
  Tạo/lifetime do OS quản lý.                 TỰ ĐỘNG PHỒNG to khi cần).
  Context switch: đắt (syscall/kernel).      Context switch: RẺ (ngay ở user space).
  Số lượng giới hạn (trăm~nghìn).            Số lượng thoải mái (trăm nghìn+).
  Scheduler: OS.                              Scheduler: runtime Go (M:N).
```

**M:N scheduling:** bạn viết `M` goroutine, runtime gán chúng vào `N` OS thread.
Bạn **không kiểm soát** thread — chỉ kiểm soát goroutine.

```text
  Goroutine A ─┐
  Goroutine B ─┼──►  Go Scheduler  ──►  OS Thread 1 (P1)
  Goroutine C ─┘         │            ──►  OS Thread 2 (P2)
   ...                   │            ──►  OS Thread 3 (P3)
                  số lượng = GOMAXPROCS
                  (= số CPU lõi, đổi bằng runtime.GOMAXPROCS)
```

---

## 2. Concurrency ≠ Parallelism

Đây là 2 khái niệm **khác nhau**, hay bị nhầm lẫn:

```text
  CONCURRENT (xen kẽ)                      PARALLEL (đồng thời thật sự)
  ──────────────────────────────────       ──────────────────────────────────
  "Nhiều việc CÙNG LÚC được QUẢN LÝ        "Nhiều việc chạy CÙNG MỘT KHOẢNH
   để tiến triển" (interleaved).            KHẮC trên nhiều lõi CPU".
  Cần 1 lõi cũng được.                      BẮT BUỘC cần >1 lõi.
  Là TÍNH CHẤT của thiết kế.               Là TÌNH TRẠNG của phần cứng.
```

```text
  1 lõi CPU, 3 goroutine xen kẽ:           2 lõi CPU, 2 goroutine song song:
  ───────────────────────────              ───────────────────────────
  CPU: [A][B][C][A][B][C]...               Lõi 1: [======A======]
  → xen kẽ, không cái nào bị "đói"        Lõi 2: [======B======]
                                            → chạy thật sự đồng thời
```

> **Goroutine cho bạn *concurrency*; còn có *parallelism* hay không** phụ thuộc
> số lõi CPU, `GOMAXPROCS` và loại công việc (CPU-bound hay IO-bound).
> Go **không tự hứa** parallelism.

```go
// ConcurrencyVsParallelism — 3 worker xử lý job xen kẽ.
// Trên máy nhiều lõi, chúng CŨNG có thể chạy song song thật sự.
func ConcurrencyVsParallelism() {
    jobs := make(chan int, 4)
    done := make(chan struct{})

    for w := 1; w <= 3; w++ {
        go func(id int) {
            for j := range jobs {
                fmt.Printf("  worker %d xử lý job %d\n", id, j)
                time.Sleep(50 * time.Millisecond) // giả lập công việc
            }
            done <- struct{}{}
        }(w)
    }

    for j := 1; j <= 4; j++ {
        jobs <- j
    }
    close(jobs)
    for w := 1; w <= 3; w++ {
        <-done
    }
}
```

Output (thứ tự worker có thể khác nhau mỗi lần chạy):

```text
  worker 1 xử lý job 1
  worker 2 xử lý job 2
  worker 3 xử lý job 3
  worker 1 xử lý job 4
```

---

## 3. Tạo goroutine — keyword `go`

Cú pháp: đặt từ khóa `go` trước **lệnh gọi hàm** là xong.

```go
func sayHello() {
    fmt.Println("Say hello from goroutine")
}

func Goroutines() {
    go sayHello()          // 1. hàm đã khai báo sẵn
    go func() {            // 2. anonymous function (closures)
        fmt.Println("lambda")
    }()

    x := 10
    go func(v int) {       // 3. anonymous function CÓ tham số —
        fmt.Println(v)     //    truyền vào tường minh, đừng capture ngoài
    }(x)

    time.Sleep(time.Second) // ← chỉ để demo, xem mục 5
}
```

```text
  Cú pháp:
  ─────────────────────────────────────────────────────────────
  go f()            — chạy hàm f trong goroutine MỚI.
  go func(){...}()  — closure; PHẢI có () ở cuối để thực sự gọi nó.
  go obj.Method()   — method cũng chạy được như hàm bình thường.
```

> ⚠️ **`go func(){...}` quên `()` ở cuối** → bạn đang truyền *tham chiếu hàm*
> cho `go`, và `go` **báo lỗi biên dịch**: *"go requires a function call"*.

### Goroutine khởi động nhanh, không đảm bảo thứ tự

```go
for i := 0; i < 3; i++ {
    go fmt.Println("in từ goroutine", i)
}
fmt.Println("in từ main")
```

```text
  KHÔNG có thứ tự nào được đảm bảo:
  ─────────────────────────────────────────────────────────────
  • Không phải 0,1,2 — có thể là 2,0,1 hay bất kỳ thứ tự nào.
  • `go` chỉ NÓP việc cho scheduler rồi trả NGAY lập tức,
    không chờ goroutine chạy xong.
```

---

## 4. Lifecycle — `main()` không chờ ai

**Quy tắc quan trọng nhất:** khi `main()` return, **chương trình kết thúc**.
Mọi goroutine đang chạy đều bị **dừng ngay lập tức** — kể cả khi chúng chưa xong.

```go
func Lifecycle() {
    go sayHello()
    // KHÔNG có dòng nào chờ → main return luôn → sayHello có thể CHƯA kịp chạy.
    fmt.Println("main xong")
}
```

Output khi **không** chờ (bản chất cuộc đua):

```text
main xong
                     ← sayHello KHÔNG in: main đã thoát trước
```

Output khi có chờ:

```go
func Lifecycle() {
    go sayHello()
    time.Sleep(time.Second) // ← KHÔNG phải cách đúng, chỉ là cách tạm
    fmt.Println("main xong")
}
```

```text
Say hello from goroutine
main xong
```

```text
  Lifecycle của 1 goroutine:
  ─────────────────────────────────────────────────────────────
  go f()  ──►  runnable (đưa vào hàng đợi scheduler)
            ──►  running (được gán 1 P, chạy trên 1 OS thread)
            ──►  blocked (đọc/ghi channel đang chờ, syscall, sleep…)
            └ main() return ──►  DỪNG TẤT CẢ, không chờ, không giải thích.
```

> **Cách đúng để chờ goroutine:** `sync.WaitGroup` hoặc **channel** như ví dụ
> ở mục 2 (`done` channel). `time.Sleep` là **phép đoán**, không phải đồng bộ
> — xem mục 6.

---

## 5. Đồng bộ sơ bộ bằng Channel

Channel là cơ chế **giao tiếp + đồng bộ** của Go: `ch <- v` gửi giá trị,
`<-ch` nhận. Unbuffered channel (không có `make(chan T, n)`) là **synchronous hand-off**:
người gửi **đứng chờ** cho tới khi có người nhận, và ngược lại.

```go
func calculate(ch chan int) {
    result := 10 + 20
    ch <- result   // gửi kết quả, CHỜ có người nhận mới chạy tiếp
}

func ChannelSync() {
    ch := make(chan int)   // unbuffered
    go calculate(ch)
    result := <-ch          // main CHỜ calculate gửi
    fmt.Println("Result from channel:", result)
}
```

Output:

```text
Result from channel: 30
```

```text
  Đồng bộ bằng unbuffered channel:
  ─────────────────────────────────────────────────────────────
  goroutine A: ch <- 30   ─┐
                           ├──►  tay bắt tay: cả 2 dừng lại 1 nhịp
  main:        v := <-ch  ─┘
  → KHÔNG cần sleep, KHÔNG cần đoán thời gian. Đây là đồng bộ THẬT.
```

> Đây là cơ chế "chờ" **chính thống** của Go: thay vì đoán bằng `Sleep`,
> để chính channel tạo điểm gặp nhau. Bài [Channels](<Channels.md>) đi sâu
> buffered / directional / close.

---

## 6. Data Race — vì sao phải cẩn thận

`counter++` **không** là 1 lệnh CPU. Nó là 3 bước: **đọc → cộng → ghi**.
2 goroutine cùng làm việc này thì kết quả bị **mất cập nhật** (lost update).

```go
var counter int

func increment() { counter++ }

func DataRace() {
    const n = 1000
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // KHÔNG an toàn: không có lock/atomic
        }()
    }
    wg.Wait()
    fmt.Printf("counter = %d (mong đợi %d)\n", counter, n)
}
```

Output **có thể** là (kết quả không ổn định):

```text
counter = 973 (mong đợi 1000)
```

```text
  Vì sao mất cập nhật:
  ─────────────────────────────────────────────────────────────
  Goroutine 1:  đọc counter(=100) ─┐
  Goroutine 2:  đọc counter(=100) ─┼─► cả 2 đều ghi 101
                                    → MỘT cập nhật bị MẤT (se lẽ phải là 102).
  Data race còn gây ra kết quả "bất thường" tùy CPU/compiler — không chỉ sai số.
```

**Phát hiện race:** chạy với cờ `-race` — đây là lý do bạn PHẢI dùng nó trong dev:

```bash
go run -race .
go test -race ./...
```

```text
  ================== DATA RACE ==================
  Write at 0x00000114d148 by goroutine 7:
    main.increment()
        fundamental/5.concurrency/goroutines.go:16 +0x18
  Previous read at 0x00000114d148 by goroutine 6:
    ...
  ===============================================
```

> **Cách sửa** (phạm vi bài sau): `sync.Mutex`, `sync/atomic`, hoặc **tránh
> chia sẻ trạng thái** bằng cách giao tiếp qua channel (*"Don't communicate
> by sharing memory; share memory by communicating"*).

---

## 7. Pitfall thường gặp

### 7.1 Lạm dụng `time.Sleep` để chờ goroutine

`Sleep` là **phép đoán thời gian**, không phải đồng bộ:

```go
go slowAPI()          // mất 2 giây, chương trình đoán 1 giây
time.Sleep(time.Second)  // → main thoát TRƯỚC khi slowAPI xong. BUG.
```

```text
  ❌  time.Sleep(d)  — d quá ngắn: mất việc. d quá dài: chương trình chậm vô ích.
  ✅  channel        — chờ đúng thời điểm có dữ liệu.
  ✅  sync.WaitGroup — chờ đúng khi các goroutine báo Done().
```

### 7.2 Panic trong goroutine làm CRASH cả chương trình

`panic` lan ra tận goroutine gốc → **không còn frame nào để recover** → process chết.
`defer + recover()` **phải** được đăng ký **trong chính goroutine đó**.

```go
func PanicInGoroutine() {
    defer func() {
        if r := recover(); r != nil {   // ← KHÔNG bắt được panic của goroutine con
            fmt.Println("main recover:", r)
        }
    }()

    go func() {
        defer func() {
            if r := recover(); r != nil { // ← MỚI là chỗ recover đúng
                fmt.Println("worker recover:", r)
            }
        }()
        panic("boom")
    }()

    time.Sleep(time.Second)
}
```

```text
  Sơ đồ panic trong goroutine:
  ─────────────────────────────────────────────────────────────
  goroutine con: panic("boom")
       │  defer trong CHÍNH goroutine đó có recover? ──► có → sống sót
       └──────────────────────────────────────────────► không → CRASH CẢ PROCESS
  → Không có "panic propagation" lên main. Đây là khác biệt lớn so với
    exception trong các ngôn ngữ khác.
```

> Xem thêm [panic and recover](<../Error Handling/panic and recover.md>).

### 7.3 Goroutine Leak — goroutine không bao giờ thoát

Goroutine bị **kẹt vĩnh viễn** ở thao tác block (channel không ai gửi/nhận,
mutex không ai mở, network không timeout) → **rò rỉ** bộ nhớ + slot scheduler.
Server chạy lâu ngày sẽ **chết dần** vì leak.

```go
func GoroutineLeak() {
    ch := make(chan int) // unbuffered

    go func() {
        ch <- 1 // block MÃI MÃI: main không còn đọc ch nữa
        fmt.Println("dòng này không bao giờ được in")
    }()

    _ = <-ch // main đọc MỘT lần rồi bỏ đi
    fmt.Println("main thoát, goroutine kia vẫn còn bị kẹt")
}
```

Output:

```text
main thoát, goroutine kia vẫn còn bị kẹt
                     ← goroutine kia vẫn nằm đó, không in gì cả
```

```text
  Dấu hiệu nhận biết leak:
  ─────────────────────────────────────────────────────────────
  • Goroutine block ở <-ch / ch <- nhưng không còn ai "đối tác".
  • Goroutine chờ mutex không bao giờ được unlock.
  • Goroutine trong vòng lặp read I/O không có timeout/cancellation.
  → Kiểm tra bằng pprof: runtime.NumGoroutine() tăng vô hạn
    hoặc `go tool pprof http://localhost:6060/debug/pprof/goroutine`.
  → Cách phòng: luôn có "lối thoát" (context.Context, done channel, close()).
```

### 7.4 Capture biến loop

Trước **Go 1.22**: biến loop là **một biến duy nhất** cho tất cả iteration →
mọi goroutine đều thấy giá trị **cuối cùng**:

```go
// BUG trên Go < 1.22 — cả 3 goroutine in i = 3
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()      // capture biến NGOÀI loop
}
```

Từ **Go 1.22+** (module của bạn đang `go 1.26.0`) → **mỗi iteration có biến mới**,
nên đoạn trên đã **an toàn**. Nhưng quy tắc truyền tham số vẫn nên giữ:

```go
// Cách tường minh, chạy đúng trên MỌI phiên bản Go
for i := 0; i < 3; i++ {
    go func(v int) { fmt.Println(v) }(i)   // truyền tường minh
}

// Hoặc, nếu capture biến ngoài loop khác (vd slice element), tự tạo bản sao:
for _, v := range items {
    v := v          // "biến mới" thủ công
    go func() { use(v) }()
}
```

```text
  Quy tắc an toàn:
  ─────────────────────────────────────────────────────────────
  • Go ≥ 1.22: biến loop đã per-iteration → closure capture i an toàn.
  • VẪN nên truyền tham số tường minh: code rõ nghĩa, không phụ thuộc
    semantics của phiên bản Go, và tránh nhầm lẫn khi đọc lại.
```

---

## 8. Best Practices / Lưu ý quan trọng

* **Tạo goroutine ở tầng cao nhất cần thiết, sở hữu nó (ownership).**
  Goroutine nào sinh ra thì phải chịu trách nhiệm cho việc kết thúc nó.
* **Luôn có lối thoát.** Mỗi goroutine cần một cách để dừng: `context.Context`,
  `done` channel, `close(ch)`, hoặc channel đóng.
* **Không dùng `time.Sleep` để đồng bộ** — dùng channel / `WaitGroup`.
* **Không để goroutine panic không có recover** — worker nên tự `recover` ở
  boundary của nó và trả `error` ra ngoài.
* **Chạy `go test -race ./...` thường xuyên** — race chỉ lộ ra khi chạy `-race`.
* **Cẩn thận truyền biến vào goroutine** — truyền tường minh qua tham số,
  đừng capture biến ngoài vòng lặp.
* **Giới hạn số lượng goroutine** — đừng `go` vô hạn; dùng worker pool
  (số goroutine hữu hạn, job đi qua channel).
* Theo dõi số goroutine đang sống bằng `runtime.NumGoroutine()` khi nghi leak.

---

## 9. Summary

* **Goroutine = lightweight task** do runtime Go scheduler quản lý (M:N),
  stack ~2-8KB, context switch rất rẻ.
* **Concurrency ≠ Parallelism:** xen kẽ là tính chất của thiết kế;
  parallelism là tình trạng của phần cứng (cần >1 lõi).
* `go f()` khởi động ngay và **không chờ**; `main()` thoát → **tất cả** goroutine bị dừng.
* **Đồng bộ bằng channel** (unbuffered = tay bắt tay) là cách chính thống —
  thay vì đoán bằng `time.Sleep`.
* **Data race** (`counter++` read-modify-write) làm mất cập nhật → phát hiện bằng
  `go run -race`, sửa bằng Mutex/atomic hoặc chia sẻ ít nhất có thể.
* **Goroutine leak** = goroutine kẹt vĩnh viễn → luôn có lối thoát (context, close, done).
* **Panic trong goroutine = crash cả process** → `defer + recover` phải nằm
  trong chính goroutine đó.

---

## Bài tiếp theo

| Bài | Nội dung |
|-----|-----------|
| [Channels](<Channels.md>) ✅ | Buffered/unbuffered, directional, close, range, **Select** |
| [Worker Pool](<Worker Pool.md>) ⬜ | `sync.WaitGroup`, pool cố định, job động, chống deadlock |
| [Synchronization](<Synchronization.md>) ⬜ | Mutex, RWMutex, Atomic |
| [Context](<Context.md>) ⬜ | Cancel, timeout, value |
| [Concurrency Patterns](<Concurrency Patterns.md>) ⬜ | Fan-Out/In, Pipeline, Producer-Consumer |
