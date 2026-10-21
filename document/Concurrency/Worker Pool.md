# Worker Pool trong Go

**Worker Pool** = một nhóm goroutine **cố định** cùng xử lý hàng đợi job.
Thay vì tạo 1 goroutine cho mỗi việc → tạo **N worker**, cho chúng ăn job.

Code mẫu: `fundamental/5.concurrency/worker_pool.go`
Đọc trước: [Goroutines](<Goroutines.md>), [Channels](<Channels.md>)

> Bài này cũng giảng **`sync.WaitGroup`** (nền tảng để chờ goroutine).

---

## 1. Tại sao cần Worker Pool?

```text
  KHÔNG pool: 1000 job → 1000 goroutine
  ─────────────────────────────────────────────────────────
  • Tạo/destroy goroutine liên tục → tốn.
  • Cạnh tranh RAM + scheduler nếu job "nặng" (DB, HTTP).
  • Không kiểm soát được độ đồng thời → DDoS chính service của bạn.

  CÓ pool: 1000 job → 50 worker (giả sử)
  ─────────────────────────────────────────────────────────
  • ĐỘNG THỜI TỐI ĐA = số worker. Job đến đâu xử lý đến đó.
  • Goroutine sống lâu, tái sử dụng.
  • Dễ kiểm soát: thêm/bớt worker, thêm timeout, thêm context.
```

```text
        jobs ──► ┌────────────┐ ──► results
        jobs ──► │ worker  x N │ ──► results
        jobs ──► └────────────┘ ──► results
        (N cố định — mỗi worker lặp lại: lấy job → làm → lặp)
```

---

## 2. `sync.WaitGroup` — chờ goroutine xong

Channel **truyền dữ liệu**; WaitGroup **đếm việc đã xong**. Đây là cách chờ
"không cần truyền gì qua channel":

```go
func WaitGroupDemo() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1) // Add TRƯỚC khi `go` — sai chỗ này là panic/race
        go func(id int) {
            defer wg.Done()              // báo "1 việc xong"
            fmt.Printf("  worker %d đang làm...\n", id)
            time.Sleep(50 * time.Millisecond)
        }(i)
    }

    wg.Wait() // block tới khi counter về 0
    fmt.Println("  tất cả worker đã xong")
}
```

Output thật:

```text
  worker 1 đang làm...
  worker 3 đang làm...
  worker 2 đang làm...
  tất cả worker đã xong
```

```text
  3 phương thức của WaitGroup:
  ─────────────────────────────────────────────────────────
  wg.Add(n)    báo TRƯỚC có n việc cần chờ   ← PHẢI chạy TRƯỚC `go`
  wg.Done()    báo 1 việc đã xong            ← defer trong goroutine
  wg.Wait()    chờ counter về 0              ← ở goroutine CHỜ (main)

  ❌  Add SAU `go`        → main có thể Wait khi counter còn 0 → không chờ gì.
  ❌  Quên Done()         → Wait block MÃI MÃI (deadlock).
  ❌  Wait khi đang Add   → Add bị gọi SAU khi Wait bắt đầu → race.
  ❌  wg dùng lại sai     → WaitGroup KHÔNG được copy sau khi đã dùng.
```

> **Chọn WaitGroup hay Channel?**
> Chỉ cần biết "khi nào xong" → `WaitGroup`.
> Cần truyền kết quả qua → truyền thêm 1 channel `results`.

---

## 3. Worker Pool cơ bản

```go
func WorkerPool() {
    const (
        numWorkers = 3
        numJobs    = 6
    )

    jobs := make(chan int, numJobs)    // buffered: nạp sẵn job
    results := make(chan int, numJobs) // buffered: không cần đọc ngay

    var wg sync.WaitGroup

    // 1. Khởi động ĐÚNG numWorkers goroutine — chúng CHẠY LẠI nhiều lần
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := range jobs { // lặp tới khi `jobs` bị close
                fmt.Printf("    worker %d xử lý job %d\n", id, j)
                result := j * j
                time.Sleep(30 * time.Millisecond) // giả lập công việc nặng
                results <- result
            }
        }(w)
    }

    // 2. Gửi toàn bộ job
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs) // báo "hết job" → worker thoát vòng `for`

    // 3. Chờ worker xong, RỒI mới đóng results
    wg.Wait()
    close(results)

    // 4. Thu kết quả
    total := 0
    for r := range results {
        total += r
    }
    fmt.Printf("  tổng kết quả = %d (%d workers, %d jobs)\n", total, numWorkers, numJobs)
}
```

Output thật (thứ tự worker **không cố định**):

```text
    worker 3 xử lý job 1
    worker 2 xử lý job 2
    worker 1 xử lý job 3
    worker 1 xử lý job 4
    worker 2 xử lý job 5
    worker 3 xử lý job 6
  tổng kết quả = 91 (3 workers, 6 jobs)
```

```text
  Thứ tự quan trọng (thứ tự SAI = deadlock):
  ─────────────────────────────────────────────────────────
  ① go worker              ──┐
  ② jobs <- job; close(jobs) ─┤  producer chạy
  ③ wg.Wait()               ─┤  CHỜ worker thoát
  ④ close(results)          ─┘  CHỈ KHI worker đã xong
  ⑤ for r := range results     thu kết quả
```

---

## 4. Job đến ĐỘNG (không biết trước số lượng)

Khi job không nạp sẵn, `results` **không đóng được ở main** — vì `wg.Wait()`
ở main sẽ chặn việc đọc results (results đầy → worker block → wg không về 0).

Cách giải quyết: `wg.Wait()` + `close(results)` đưa vào **1 goroutine riêng**:

```go
func WorkerPoolDynamic() {
    jobs := make(chan int)
    results := make(chan int)

    var wg sync.WaitGroup
    for w := 1; w <= 2; w++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := range jobs {
                results <- j * 10
            }
        }(w)
    }

    // producer: tự quyết định lúc nào gửi, lúc nào đóng
    go func() {
        for _, v := range []int{5, 8, 13, 21} {
            jobs <- v
        }
        close(jobs)
    }()

    // Đóng results trong goroutine RIÊNG: wg.Wait() không phải main
    go func() {
        wg.Wait()
        close(results)
    }()

    sum := 0
    for r := range results {
        sum += r
    }
    fmt.Println("  tổng kết quả (dynamic) =", sum) // 470
}
```

Output thật:

```text
  tổng kết quả (dynamic) = 470
```

```text
  Pattern "close results từ goroutine riêng":
  ─────────────────────────────────────────────────────────
  jobs ──► [worker x N] ──► results ──► main (range results)
                     │
  wg.Wait() + close(results) chạy SONG SONG với main range.
  → main không bao giờ chờ wg, wg không bao giờ chờ main → không deadlock.
```

---

## 5. Lỗi thường gặp

| Hành vi | Kết quả |
|---------|---------|
| `wg.Add(1)` **sau** `go` | main `Wait` khi counter còn 0 → **không chờ gì** |
| Quên `wg.Done()` | `wg.Wait()` **block mãi** → deadlock |
| `close(results)` **trước** `wg.Wait()` | Worker còn đang gửi → **panic**: `send on closed channel` |
| `results` buffered **quá nhỏ** + main chưa đọc | Worker gửi → **block** → wg không về 0 → deadlock |
| Tạo goroutine **không giới hạn** (`go` mỗi job) | 1000 job → 1000 goroutine → tốn RAM, cạnh tranh |
| `jobs` **không close** | Worker `range jobs` **block mãi** → leak |
| Copy `wg` sau khi đã dùng | **race** / panic: `WaitGroup is reused before previous Wait has returned` |

```text
  Cách phát hiện:
  ─────────────────────────────────────────────────────────
  • go test -race ./...        ← race + deadlock trong test
  • runtime.NumGoroutine()     ← tăng vô hạn = leak
  • go tool pprof /debug/pprof/goroutine  ← xem goroutine đang block ở đâu
```

---

## 6. Best Practices

* **Số worker = số công việc "độ đồng thời" muốn có** — thường = `GOMAXPROCS`
  (nhiều CPU-bound) hoặc lớn hơn nhiều (IO-bound: DB/HTTP chờ mạng).
* **Luôn đóng `jobs`** ở producer, **luôn `defer wg.Done()`** trong worker.
* **`close(results)` sau `wg.Wait()`** — hoặc đưa 2 việc đó vào goroutine riêng.
* **Không truyền gì, chỉ chờ** → `WaitGroup`; **cần kết quả** → thêm channel.
* Nếu job dài (HTTP, DB) → thêm **`context.Context`** để hủy pool khi shutdown
  (bài [Context](<Context.md>)).
* Đừng `go` mỗi job — đó chính là thứ Worker Pool thay thế.

---

## 7. Summary

* **Worker Pool** = N goroutine cố định ăn job từ channel → kiểm soát độ đồng thời.
* **`WaitGroup`**: `Add` **trước** `go`, `defer Done()`, `Wait()` ở nơi chờ.
* Thứ tự chuẩn: `go worker` → `close(jobs)` → `wg.Wait()` → `close(results)` → `range results`.
* Job động → `wg.Wait()` + `close(results)` chạy trong **goroutine riêng**.
* Lỗi kinh điển: quên `Done`, close sai thứ tự, buffered quá nhỏ, không đóng `jobs`.

---

## Bài tiếp theo

| Bài | Nội dung |
|-----|-----------|
| [Concurrency Patterns](<Concurrency Patterns.md>) ⬜ | Fan-Out/Fan-In, Pipeline, Producer-Consumer |
| [Context](<Context.md>) ⬜ | Hủy pool khi shutdown, timeout |
| [Synchronization](<Synchronization.md>) ⬜ | Mutex, RWMutex, Atomic (chia sẻ trạng thái) |
