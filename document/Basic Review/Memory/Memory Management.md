# Memory Management trong Go (Golang)

Go tự quản lý bộ nhớ — developer không cần `free`, **Garbage Collector** dọn dẹp.
Hiểu **stack/heap** và **escape analysis** giúp giải thích "biến nào nằm đâu" và vì sao code nhanh/chậm.

Khớp node **Memory Model** của roadmap.sh (§10) + mục **Garbage Collector** (§12).
Code mẫu: `fundamental/3.memory/memory.go`

---

## 1. Stack vs Heap

### Stack (ngăn xếp)

* Mỗi goroutine có 1 stack riêng — lưu **function call frame**: tham số, biến local, return address
* Frame bị **xóa ngay** khi hàm return → cực nhanh, không cần GC
* Biến nằm stack sẽ **không** được GC theo dõi

### Heap (pile)

* Vùng bộ nhớ do **GC quản lý** — biến ở đây sống lâu hơn hàm tạo ra nó
* Truy cập chậm hơn stack (phải qua con trỏ + GC track)
* Biến nào "trốn thoát" khỏi frame → compiler đẩy lên đây

### Sơ đồ: biến escape từ stack sang heap

```text
Stack (của goroutine)                 Heap (GC quản lý)
┌──────────────────────┐
│ frame: main           │             ┌─────────────┐
│   p ─────────────────────────────────▶  x = 42   │
├──────────────────────┤             │    (int)    │
│ frame: escapeToHeap   │             └─────────────┘
│   x  (đã chuyển lên) ─┘
└──────────────────────┘
```

### Biến nằm đâu — do **hành vi**, không do kiểu dữ liệu

Cùng là kiểu `int` trong code mẫu:

```go
func stayInStack() {
    x := 42
    _ = x            // x chỉ dùng cục bộ → STACK
}

func escapeToHeap() *int {
    x := 42
    return &x        // return địa chỉ của x → HEAP
}
```

Go **không** có quy tắc "kiểu X nằm ở heap" — tất cả do escape analysis quyết định (xem mục 2).

### Variable lifetime

| | Stack | Heap |
|---|---|---|
| Sống đến khi nào? | Hàm return là xóa | Không còn reference nào → GC thu hồi |
| Chi phí | Rất thấp (chỉ di chuyển stack pointer) | Có chi phí GC + indirection |

---

## 2. Escape Analysis

Compiler phân tích lúc **build**: biến nào "thoát" khỏi frame → đẩy lên heap.
Muốn biết biến nào escape, dùng flag — **đừng đoán**:

### Xem kết quả

```bash
go build -gcflags="-m" ./fundamental/3.memory/
```

Output thật (lọc ra file `memory.go`):

```text
fundamental\3.memory/memory.go:8:2: moved to heap: x
fundamental\3.memory/memory.go:20:2: moved to heap: n
fundamental\3.memory/memory.go:21:9: func literal escapes to heap
fundamental\3.memory/memory.go:29:14: 42 escapes to heap
```

> Chi tiết flag & toolchain: xem [Commands](<../../Tooling/Commands.md>)

### 3 lý do escape điển hình (đều có trong code mẫu)

#### 1. Return biến local dưới dạng pointer

```go
func escapeToHeap() *int {
    x := 42
    return &x        // ❌ line 8:  moved to heap: x
}
```

Hàm đã return nhưng `x` vẫn được trỏ tới → `x` phải sống trên heap.

#### 2. Biến bị closure capture

```go
func closureCapture() func() int {
    n := 0
    return func() int {   // ❌ line 20: moved to heap: n
        n++
        return n
    }
}
```

Closure có thể được gọi **sau khi** hàm ngoài đã return → `n` phải lên heap.

#### 3. Boxing vào interface

```go
func boxing() {
    var v any = 42       // ❌ line 29: 42 escapes to heap
    fmt.Println(v)
}
```

Giá trị `int` bị "đóng gói" vào `any` (interface) → compiler không biết trước concrete type → copy lên heap.

### Không escape thì sao

```go
func stayInStack() {
    x := 42
    _ = x            // ✅ KHÔNG có dòng "moved to heap" nào cho biến ở line 15
}
```

Chạy thử chương trình:

```bash
go run ./fundamental/3.memory/
```

Output thật:

```text
escapeToHeap: 42
closureCapture: 1 2
42
```

### Mẹo thực tế

* Escape **không phải lúc nào cũng xấu** — dữ liệu sống lâu hơn hàm thì *bắt buộc* phải escape
* Đo allocation thật bằng benchmark: `go test -benchmem` (xem [Commands](<../../Tooling/Commands.md>))
* "Zero-allocation" là mục tiêu khi viết hot path, không phải mọi nơi

---

## 3. Value Semantics (tóm tắt)

Go **luôn** pass by value — gán biến hay truyền tham số đều **copy**:

| | Value (`T`) | Pointer (`*T`) |
|---|---|---|
| Cái được copy | Toàn bộ dữ liệu | Chỉ địa chỉ (8 byte) |
| Sửa được bản gốc? | ❌ Không | ✅ Có |
| Ảnh hưởng memory | Copy dữ liệu (tốn nếu struct lớn) | Copy địa chỉ (rẻ) |

```go
a := User{Name: "Trung"}   // copy toàn bộ field
b := &a                     // chỉ copy địa chỉ — a và b dùng chung dữ liệu
```

> Phân tích kỹ hơn ở [Pointer](<../../Composite Types/Pointer.md>) — mục *Pass by value vs pointer*,
> và [Struct](<../../Composite Types/Struct.md>).

**Gợi ý nhanh:** struct nhỏ, hàm chỉ đọc → truyền value.
Cần mutate gốc, struct lớn, hoặc receiver có method mutate → truyền pointer.

---

## 4. Garbage Collector (tóm tắt)

Go dùng **tracing GC, concurrent Mark & Sweep** — developer không gọi `free`,
mọi object trên heap do GC thu hồi. Nhờ **escape analysis**, compiler còn biết biến nào nằm stack (không cần GC theo dõi) hay heap (cần GC).

Quan sát nhanh:

```bash
$env:GODEBUG="gctrace=1"; go run ./fundamental/3.memory/
```

> **Chi tiết đầy đủ** (tri-color mark & sweep, write barrier, giải nghĩa dòng `gc 1 ...`,
> `GOGC` vs `GOMEMLIMIT`, `runtime.ReadMemStats`, finalizer):
> xem [Garbage Collector](<Garbage Collector.md>)

---

## 5. Summary

* **Stack**: frame của hàm, biến local không escape, xóa ngay khi return — cực nhanh
* **Heap**: vùng GC quản lý — biến sống lâu hơn hàm tạo ra nó
* Biến nằm ở đâu do **escape analysis** quyết định, **không** do kiểu dữ liệu
* `go build -gcflags="-m"` để thấy `moved to heap:` — đừng đoán
* 3 lý do escape phổ biến: **return pointer**, **closure capture**, **boxing vào interface**
* Go GC: **Mark & Sweep** chạy song song — chi tiết ở [Garbage Collector](<Garbage Collector.md>)
* Escape không phải lúc nào cũng xấu — chỉ tối giảm allocation ở hot path
