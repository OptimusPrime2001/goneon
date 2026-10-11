# Map trong Go (Golang)

`map` là cấu trúc dữ liệu **key → value** kiểu tra cứu nhanh (hash table), là **reference type**:
biến map chỉ giữ con trỏ tới bảng hash, nên gán / truyền qua hàm đều **dùng chung dữ liệu gốc**.

> Thuật ngữ xuất hiện trong roadmap.sh → mục **Composite Types → Maps → Comma-Ok Idiom**.
> Code mẫu: `fundamental/2.composite_type/map.go` (đang tổng hợp từ `fundamental/3.reference_types/map.go`).

---

## 1. Khái niệm

```go
var m map[string]int          // nil map
m["Apple"] = 20               // ❌ panic: assignment to entry in nil map
```

| Đặc điểm | Giải thích |
|-----------|-----------|
| Key → Value | Mỗi key trỏ đúng 1 value, key **không được lặp lại** |
| Dynamic | Tự tăng kích thước, không cần khai báo độ dài |
| Reference type | Gán `b := a` → `a` và `b` **dùng chung** 1 bảng hash |
| Không an toàn concurrency | Nhiều goroutine ghi/đọc cùng lúc phải tự khóa (xem mục 12) |
| Zero value | Là `nil` — đọc được, **ghi panic** |

---

## 2. Khởi tạo

### Cách 1: literal có sẵn phần tử

```go
giaTraiCay := map[string]int{
    "Apple":  20,
    "Orange": 25,
    "Banana": 31,
}
```

### Cách 2: `make` — **bắt buộc** nếu muốn ghi vào map

```go
m := make(map[string]int)          // rỗng
m := make(map[string]int, 1000)    // capacity hint: dự báo số phần tử, giảm lần cấp phát lại
```

> ⚠️ `make(map[string]int, 1000)` là **gợi ý dung lượng**, không phải độ dài cố định — map vẫn tăng thêm bình thường.
> Không dùng được `cap(m)` trên map:

```text
invalid argument: big (variable of type map[string]int) for built-in cap
```

(`cap` chỉ nhận array, *array, slice, channel.)

### Cách 3: literal rỗng

```go
m := map[string]int{}   // KHÔNG phải nil — ghi vào bình thường
```

### Cách 4: khai báo biến rồi gán

```go
var m map[string]int    // nil
m = map[string]int{}    // vẫn nil -> chưa ghi được
m = make(map[string]int) // OK
```

### So sánh 3 trạng thái

| | `var m map[string]int` | `m := map[string]int{}` | `m := make(map[string]int)` |
|---|---|---|---|
| Là nil? | ✅ có | ❌ không | ❌ không |
| Đọc `m["x"]` | trả zero value | trả zero value | trả zero value |
| Ghi `m["x"]=1` | **panic** | OK | OK |
| `len(m)` | 0 | 0 | 0 |

---

## 3. Thêm / Sửa / Cập nhật

```go
giaTraiCay := map[string]int{"Apple": 20}

giaTraiCay["Apple"] = 22      // key đã có  -> cập nhật
giaTraiCay["Mango"] = 26      // key chưa có -> thêm mới
```

* Một cú pháp duy nhất cho cả "thêm" và "cập nhật" — Go **không** phân biệt.
* Key không phân biệt HOA/thường: `"apple"` ≠ `"Apple"` là 2 key khác nhau.
* Key `""` (chuỗi rỗng) hoặc `0` là key hợp lệ bình thường.

---

## 4. Đọc giá trị — Comma-Ok Idiom ⭐

### Cách 1: truy cập trực tiếp

```go
gia := giaTraiCay["Apple"]
```

Nếu key **không tồn tại** → Go trả **zero value của value type** (`0`, `""`, `false`, `nil`), **không panic**.

> ⚠️ **Với map `bool`, cách này KHÔNG phân biệt được "key không có" và "value là false"**:

```go
online := map[string]bool{"a": true}
fmt.Println(online["b"])       // false  — key "b" không hề tồn tại
```

### Cách 2: Comma-Ok Idiom — cách đúng

```go
gia, exists := giaTraiCay["Apple"]
if exists {
    fmt.Println("Apple price is", gia)
} else {
    fmt.Println("Apple chưa có trong map")
}
```

Output thật:

```text
Apple: 20 true
Grape: 0 false     <- "Grape" không tồn tại, gia = 0 là zero value
```

| Biến | Ý nghĩa |
|------|---------|
| `v, ok := m[k]` | `ok == true` → key có thật, `v` là giá trị |
| | `ok == false` → key **không có**, `v` = zero value |

Dùng cho cả map rỗng và nil map:

```go
var m map[string]int
v, ok := m["x"]   // v = 0, ok = false  → KHÔNG panic
```

---

## 5. Xoá phần tử

```go
delete(giaTraiCay, "Apple")   // xoá key "Apple"
delete(giaTraiCay, "Cam")     // key không tồn tại -> no-op, KHÔNG panic
fmt.Println(len(giaTraiCay))  // số phần tử hiện tại
```

* `delete` trên **nil map** cũng không panic (no-op).
* Go 1.21+ có builtin `clear(m)` để xoá **tất cả** cùng lúc:

```go
clear(giaTraiCay)
fmt.Println(giaTraiCay, len(giaTraiCay))   // map[] 0
```

---

## 6. Duyệt map — `range`

```go
for k, v := range giaTraiCay {
    fmt.Printf("Qua %s co gia la %d\n", k, v)
}
```

| Cú pháp | Lấy được |
|---------|----------|
| `for k, v := range m` | key + value |
| `for k := range m` | **chỉ key** (nhanh hơn) |
| `for range m` | không lấy gì, chỉ duyệt cho hết |
| `for _, v := range m` | **chỉ value** |

### 3 điều quan trọng

1. **Thứ tự là NGẪU NHIÊN** — Go spec nói *"the iteration order over maps is not specified"*. Muốn có thứ tự phải tự sắp xếp key:

```go
keys := make([]string, 0, len(giaTraiCay))
for k := range giaTraiCay {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, giaTraiCay[k])
}
```

2. **Mỗi key được duyệt đúng 1 lần** — dù nội dung thay đổi giữa chừng, không duyệt trùng.

3. **Xóa key đang duyệt thì an toàn** (đây là pattern dọn map phổ biến):

```go
for k := range m {
    delete(m, k)     // ✅ an toàn
}
// kết quả: map[] len: 0
```

Nhưng **thêm key mới trong lúc range** thì key đó có thể được duyệt hoặc không — **không được dựa vào**.

---

## 7. Reference semantics — gán = dùng chung ⭐

```go
a := map[string]int{"x": 1}
b := a          // KHÔNG copy bảng hash, chỉ copy con trỏ
b["x"] = 99

fmt.Println("a =", a)   // a = map[x:99]   ← bị ảnh hưởng!
fmt.Println("b =", b)   // b = map[x:99]
```

**Truyền vào hàm cũng dùng chung:**

```go
func sua(m map[string]int) {
    m["x"] = 999
}

func main() {
    a := map[string]int{"x": 1}
    sua(a)
    fmt.Println(a)      // map[x:999]
}
```

### ⚠️ Bẫy kinh điển: copy struct **vẫn dùng chung map** bên trong

```go
type SV struct {
    Ten string
    Di  map[string]int
}

s1 := SV{Ten: "HN", Di: map[string]int{"HN": 1}}
s2 := s1       // copy struct: Ten bị copy, NHƯNG Di vẫn là map tham chiếu

s2.Ten = "HCM"
s2.Di["HN"] = 99

fmt.Printf("s1=%+v s2=%+v\n", s1, s2)
```

Output thật:

```text
s1={Ten:HN Di:map[HN:99]} s2={Ten:HCM Di:map[HN:99]}
```

→ `Ten` tách rời, nhưng `Di` bị dính. Đây là lỗi rất khó tìm nếu struct chứa map/slice.

---

## 8. Deep copy — phải tự tay

Go **không** có hàm copy map built-in. Phải copy thủ công:

```go
src := map[string]int{"a": 1, "b": 2}
dst := make(map[string]int, len(src))
for k, v := range src {
    dst[k] = v
}
dst["a"] = 100

fmt.Println("src =", src)   // src = map[a:1 b:2]    ← giữ nguyên
fmt.Println("dst =", dst)   // dst = map[a:100 b:2]
```

> ⚠️ Cách trên là **shallow copy**: nếu value là map/slice thì copy vẫn tham chiếu.

---

## 9. Kiểu key hợp lệ — `comparable`

Key phải là kiểu **so sánh được**:

| Dùng được | ❌ KHÔNG dùng được |
|-----------|--------------------|
| `int`, `float`, `bool`, `string` | `slice` — `map[[]int]string{}` → compile error |
| `pointer` | `map` — `map[map[string]int]string{}` |
| `array` (tất cả phần tử comparable) | `func` |
| `struct` (tất cả field comparable) | |
| `interface` (nếu dynamic type comparable) | |

```text
invalid map key type []int
```

Struct làm key rất hữu ích (vd map tọa độ):

```go
type P struct{ X, Y int }
pairs := map[P]string{{1, 2}: "a"}
fmt.Println(pairs)   // map[{1 2}:a]
```

---

## 10. Map lồng nhau & tổ hợp

### Map của slice

```go
m := map[string][]int{
    "so": {1, 2, 3},
    "le": {1, 3},
}
```

### Map lồng nhau — nhớ `make` tầng trong

```go
nested := map[string]map[string]int{}
nested["VN"] = map[string]int{}     // BẮT BUỘC, nếu không -> panic khi ghi
nested["VN"]["HN"] = 1
```

Hoặc gán trực tiếp:

```go
nested["VN"]["HCM"] = 2      // ❌ panic nếu nested["VN"] chưa tồn tại
nested["LA"] = map[string]int{"VTE": 3}   // ✅
```

> 💡 **Đọc** map lồng khi tầng trong là nil thì **không** panic (trả zero value),
> chỉ **ghi** mới panic. Nên `fmt.Println(nested["VN"]["HN"])` vẫn chạy được.

Output thật:

```text
map[LA:map[VTE:3] VN:map[HCM:2 HN:1]]
```

### Map của struct

```go
type SV struct {
    Ten string
    Di  map[string]int
}
ds := map[string]SV{
    "HN": {Ten: "Ha Noi", Di: map[string]int{"HN": 1}},
}
```

---

## 11. Nil map — đọc được, ghi panic

```go
var m map[string]int

fmt.Println(len(m))        // 0        ✅
fmt.Println(m["x"])        // 0        ✅  (zero value)
v, ok := m["x"]            // 0 false  ✅
delete(m, "x")             // no-op    ✅
m["x"] = 1                 // ❌ panic
```

Output thật:

```text
len: 0
doc (không panic): 0
comma-ok: 0 false
delete nil map: OK
ghi vào nil map -> assignment to entry in nil map
```

**Bài học:** khai báo `var m map[string]int` rồi truyền vào hàm là **an toàn** (chỉ đọc).
Muốn hàm ghi được thì hàm phải `make` hoặc nhận map đã `make`.

---

## 12. An toàn concurrency

`map` **không** an toàn khi nhiều goroutine đọc/ghi cùng lúc → chạy với `-race` sẽ báo:

```text
fatal error: concurrent map read and map write
```

Hai cách xử lý:

```go
// Cách 1: sync.RWMutex — đọc nhiều, ghi ít
var mu sync.RWMutex
mu.RLock();  v := m[k]; mu.RUnlock()   // đọc
mu.Lock();   m[k] = v;   mu.Unlock()   // ghi

// Cách 2: sync.Map — tối ưu cho trường hợp ghi ít, key ổn định
var sm sync.Map
sm.Store("a", 1)
v, ok := sm.Load("a")
```

---

## 13. So sánh nhanh: Map vs Slice vs Struct

| | **Map** | **Slice** | **Struct** |
|---|---|---|---|
| Truy cập theo | **key** | **chỉ số** `m[0]` | **tên field** `s.Ten` |
| Đảo key/value | ✅ | ❌ | ❌ |
| Số phần tử | thay đổi | thay đổi | **cố định** (khai báo sẵn) |
| Key hợp lệ | comparable | index `int` | tên field, không đổi |
| Thứ tự | ❌ ngẫu nhiên | ✅ theo chỉ số | theo thứ tự khai báo |
| So sánh `==` | ❌ không so được | ❌ (trừ nil) | ✅ nếu mọi field comparable |
| Type | `map[K]V` | `[]T` | `struct{...}` |
| Reference/value | **reference** | reference-ish (header) | **value** |

---

## 14. Lỗi thường gặp

| Lỗi | Nguyên nhân | Cách sửa |
|-----|-------------|----------|
| `assignment to entry in nil map` | Ghi vào map chưa `make` | `m = make(map[string]int)` hoặc dùng literal `{}` |
| `invalid map key type []int` | Key là slice/map/func (không comparable) | Dùng kiểu comparable (int, string, struct…) |
| `invalid argument ... for built-in cap` | Dùng `cap(m)` trên map | Không có `cap(map)` — dùng `len` |
| Đọc ra `false` nhưng key có thật? | Map `bool` + truy cập trực tiếp | Dùng **comma-ok**: `v, ok := m[k]` |
| Sửa map gốc khi chỉ muốn copy | Map là reference type | Copy thủ công bằng `for range` |
| `fatal error: concurrent map writes` | Nhiều goroutine ghi một map | `sync.RWMutex` hoặc `sync.Map` |
| `panic` khi `nested["a"]["b"] = 1` | Chưa khởi tạo map tầng trong | `nested["a"] = map[string]int{}` trước |
| Key bị "biến mất" | Range thứ tự ngẫu nhiên + mong đợi có thứ tự | Tự sort key trước khi in |
| Muốn map so sánh `==` được | Go không cho map comparable | Dùng `maps.Equal` (Go 1.21+, `maps` package) |

---

## 15. Cheat sheet

| Việc cần làm | Lệnh |
|--------------|------|
| Tạo map rỗng | `m := make(map[string]int)` / `m := map[string]int{}` |
| Thêm / sửa | `m["k"] = v` |
| Đọc an toàn | `v, ok := m["k"]` |
| Xoá 1 key | `delete(m, "k")` |
| Xoá hết | `clear(m)` *(Go 1.21+)* |
| Đếm phần tử | `len(m)` |
| Duyệt | `for k, v := range m { }` |
| Duyệt có thứ tự | sort key slice rồi `for _, k := range keys` |
| Copy độc lập | `for k, v := range src { dst[k] = v }` |
| So sánh 2 map | `maps.Equal(a, b)` *(Go 1.21+)* |
| An toàn concurrency | `sync.RWMutex` / `sync.Map` |

**Builtin liên quan:** `len`, `delete`, `clear` — map **không** có `cap`, `append`, `copy`.

---

## 16. Bài tập

* Đếm tần suất xuất hiện của mỗi từ trong một slice string → `map[string]int`.
* Đảo key/value của map (bắt buộc kiểm tra key trùng).
* Viết hàm `deepCopy(m map[string][]int) map[string][]int`.
* Xây `map[string]map[string]int` biểu diễn điểm số `lớp → học sinh → điểm`,
  in ra top 3 bạn điểm cao nhất (cần sort).
* Viết test table-driven cho hàm tra cứu map với case: key có / key không có / nil map.

---

**Tham khảo:** [Map types — Go spec](https://go.dev/ref/spec#Map_types) · [Go maps in depth](https://go.dev/blog/maps) · [Go maps introduction](https://go.dev/blog/go-maps-introduction)
