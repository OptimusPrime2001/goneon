# Garbage Collector (GC) trong Go (Golang)

Go tự quản lý bộ nhớ: developer **không bao giờ** gọi `free` — mọi object trên heap do **GC** thu hồi.

Khớp node **Garbage Collector** (§12 Runtime Internals) của roadmap.sh.
Code mẫu: `fundamental/3.memory/gc.go`
Doc liên quan: [Memory Management](<Memory Management.md>) — stack/heap/escape analysis.

---

## 1. Tổng quan — GC của Go

| Đặc điểm | Ý nghĩa |
|---|---|
| **Tracing GC** | Quét theo tham chiếu (không reference counting như Python/Java) |
| **Concurrent** | Phần lớn chạy **song song** với chương trình — không "đứng máy" |
| **Tri-color mark & sweep** | Thuật toán đánh dấu 3 màu rồi quét (mục 2) |
| **Non-generational** | Go 1.x **không** chia young/old generation (khác Java) |
| **Design goal** | Pause < **1 ms** kể cả với heap lớn |

---

## 2. Nguyên lý — Tri-color Mark & Sweep

### 3 màu

| Màu | Nghĩa |
|---|---|
| **White** | Chưa được GC chạm tới — cuối chu kỳ còn trắng = **rác** → bị giải phóng |
| **Grey** | Đã được mark nhưng còn con chưa quét (đang chờ duyệt) |
| **Black** | Đã mark + con cái quét xong = **sống** |

### 1 chu kỳ GC gồm 4 bước

1. **Mark setup** — STW (stop-the-world) rất ngắn: bật *write barrier*, quét các biến trên stack (roots)
2. **Concurrent mark** — quét từ roots (stack, globals) theo kiểu BFS: white → grey → black; chương trình vẫn chạy
3. **Mark termination** — STW ngắn lần 2: tắt write barrier, dọn dẹp cuối
4. **Concurrent sweep** — trả vùng nhớ của các object white về allocator, cũng chạy song song

```text
roots (stack variables, globals, registers)
        │  mark
        ▼
      GREY ──── quét con ────▶ BLACK   (sống, không bị giải phóng)
        │
        │  hết chu kỳ mà vẫn WHITE
        ▼
    bị sweep — trả bộ nhớ về allocator
```

### GC assist

Nếu chương trình **allocate nhanh** trong lúc GC đang mark, Go bắt goroutine đó "giúp" GC (**GC assist**) để tiến độ mark theo kịp rate tạo rác.
→ Đây là lý do "nhiều allocation = nhiều GC assist = chậm" — tối giảm allocation ở hot path.

---

## 3. Write Barrier

### Vấn đề

Mark chạy **song song** với chương trình (mutator). Nếu mutator ghi:

```text
A.child = C     (trong lúc GC đã mark xong A nhưng chưa quét C)
```

GC có thể tưởng `C` là rác và giải phóng — trong khi `A` vẫn đang trỏ tới.

### Giải pháp của Go

Go dùng **hybrid write barrier** (từ Go 1.8): mọi thao tác **ghi pointer** đều "báo" cho GC — object đang được trỏ tới sẽ **không** bị coi là rác.

Nhờ barrier:

* GC không cần "re-scan" toàn bộ stack/heap ở cuối chu kỳ
* 2 lần STW mỗi chu kỳ chỉ mất vài **mili giây**

---

## 4. Quan sát GC

### 4.1 `GODEBUG=gctrace=1` — theo dõi từng chu kỳ

```bash
# PowerShell
$env:GODEBUG="gctrace=1"
go run ./fundamental/3.memory/
```

Output thật (dòng GC đầu tiên):

```text
gc 1 @0.046s 2%: 0+5.3+0 ms clock, 0+2.8/10/0+0 ms cpu, 4->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P
```

Giải nghĩa từng phần:

| Phần | Nghĩa |
|---|---|
| `gc 1` | GC cycle thứ 1 |
| `@0.046s` | Thời điểm chạy (tính từ lúc program start) |
| `2%` | Tỷ lệ % CPU mà GC chiếm |
| `0+5.3+0 ms` | STW mark-terminate + concurrent mark + STW mark-terminate |
| `4->4->1 MB` | Heap: đầu GC → cuối GC → **live** (sau sweep) |
| `4 MB goal` | Ngưỡng heap kích hoạt GC kế tiếp (≈ live × (1 + GOGC/100)) |
| `0 MB stacks` | Dung lượng stack |
| `12 P` | Số *logical processor* (P) mà scheduler dùng |

> Số dòng `gc` tỷ lệ thuận với **allocation rate** — chương trình nhỏ cũng thấy vài chục dòng vì chính runtime Go cũng allocate.

### 4.2 `runtime.ReadMemStats` — đọc số liệu trong code

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)

fmt.Println("HeapAlloc:", m.HeapAlloc/1024, "KB")
fmt.Println("NumGC:", m.NumGC)
```

Output thật:

```text
HeapAlloc: 173 KB
NumGC: 0
```

### 4.3 `runtime.GC()` — ép GC chạy ngay

```go
runtime.GC()
runtime.ReadMemStats(&m)
fmt.Println("NumGC sau runtime.GC():", m.NumGC)
```

Output thật:

```text
NumGC sau runtime.GC(): 1
```

> Dùng khi **benchmark/test** để đo đồng nhất — không dùng trong code production (mọi lần gọi =1 chu kỳ GC đầy đủ, tốn).

---

## 5. Điều khiển GC — `GOGC` & `GOMEMLIMIT`

### `GOGC` — ngưỡng kích hoạt GC

| Giá trị | Ý nghĩa |
|---|---|
| `GOGC=100` (mặc định) | Heap live data tăng 100% so với lần GC trước → chạy GC |
| `GOGC=200` | Ít GC hơn, tốn RAM hơn |
| `GOGC=off` | Tắt kích hoạt theo % — dùng khi đã set `GOMEMLIMIT` hoặc khi benchmark |

```bash
$env:GOGC="200"; go run ./fundamental/3.memory/
```

### `GOMEMLIMIT` — soft memory limit (Go 1.19+)

`GOGC`theo **%** — khó định lượng trong container. `GOMEMLIMIT` đặt **giới hạn mềm tính bằng bytes**:

* GC sẽ nỗ lực giữ tổng bộ nhớ heap **dưới limit**
* Vượt nhẹ → GC chạy **dồn dập** để giải phóng (nhưng **không** crash — limit là "mềm")
* Kết hợp `GOGC=off GOMEMLIMIT=1GiB` = chế độ "chỉ theo limit" — chống bị OOM-kill trong Docker/K8s

Trong code:

```go
debug.SetMemoryLimit(64 * 1024 * 1024) // 64 MiB
limit := debug.SetMemoryLimit(-1)      // -1 = đọc limit hiện tại, không đổi
fmt.Println("MemoryLimit:", limit/1024/1024, "MB")
```

Output thật:

```text
MemoryLimit: 64 MB
```

---

## 6. Finalizer — `runtime.SetFinalizer`

Finalizer cho phép **đăng ký hàm chạy khi object thành rác**:

```go
obj := new(int)
done := make(chan struct{})

runtime.SetFinalizer(obj, func(*int) {
    close(done)
})

obj = nil          // bỏ reference cuối cùng tới obj
runtime.GC()       // chu kỳ 1: phát hiện obj là rác, đưa vào finalizer queue
runtime.GC()       // chu kỳ 2: finalizer goroutine thường chạy giữa 2 chu kỳ

select {
case <-done:
    fmt.Println("finalizer: đã chạy sau GC")
case <-time.After(100 * time.Millisecond):
    fmt.Println("finalizer: CHƯA chạy — thời điểm không được guarantee")
}
```

Output thật:

```text
finalizer: đã chạy sau GC
```

### ⚠️ Caveats quan trọng

* Finalizer chạy trên **1 goroutine riêng, ASYNC** — thời điểm **không được guarantee** (thậm chí không chạy trước khi program thoát)
* Cần ít nhất **2 chu kỳ GC** để finalizer được xử lý
* **Không dùng finalizer để** đóng file/release resource critical → dùng `defer` + `Close()` tường minh
* Object có finalizer được GC giữ thêm1 chu kỳ → tốn bộ nhớ lâu hơn

---

## 7. Summary

* Go dùng **tracing GC, concurrent, tri-color mark & sweep** — không reference counting, không generational
* 1 chu kỳ: **mark setup (STW) → concurrent mark → mark termination (STW) → concurrent sweep** — 2 lần STW chỉ vài ms
* **Write barrier** (hybrid, Go 1.8+) bảo đảm concurrent mark không bỏ sót object
* **GC assist**: allocate nhanh trong lúc GC chạy = bị ép "giúp" GC
* Quan sát: `GODEBUG=gctrace=1` (đọc `3->4->1 MB` = đầu GC → cuối GC → live), `runtime.ReadMemStats`, ép chạy bằng `runtime.GC()`
* `GOGC=100` mặc định; **`GOMEMLIMIT`** (Go 1.19+) đặt giới hạn mềm — rất hữu ích trong container (`GOGC=off GOMEMLIMIT=1GiB`)
* **Finalizer** (`SetFinalizer`) chạy async, không guarantee — chỉ là "tính năng phụ", không phải cơ chế release resource
