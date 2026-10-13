# Generics trong Go (Golang)

Generics (Go 1.18+) cho phép viết **1 hàm/kiểu dùng được cho nhiều kiểu dữ liệu** — giữ nguyên type safety.

Khớp node **Generics** của roadmap.sh (§9).
Code mẫu: `fundamental/4.method_interface/generics.go`

---

## 1. Why — tại sao cần generics?

Trước Go 1.18, mỗi kiểu cần **1 hàm riêng**:

```go
func SumInt(nums ...int) int {            // bản cho int
    var total int
    for _, n := range nums { total += n }
    return total
}

func SumFloat(nums ...float64) float64 {  // copy y hệt, chỉ đổi kiểu
    var total float64
    for _, n := range nums { total += n }
    return total
}
```

Cách thứ 2 ngày xưa: `interface{}` — nhưng **mất type safety**:

```go
func SumAny(nums ...interface{}) int {
    total := 0
    for _, n := range nums {
        total += n.(int)      // mỗi phần tử phải assert — sai type = PANIC lúc runtime
    }
    return total
}
```

→ Generics giải quyết cả hai: **1 hàm** + **kiểm tra kiểu lúc compile**.

---

## 2. Type Parameters

```go
func Sum[T Number](nums ...T) T { ... }
//     ↑↑↑↑↑↑↑↑↑
//     [T Number] = khai báo type parameter T với constraint Number
```

* `T` là **biến tạm cho kiểu** — thay bằng int, float64, string… khi gọi
* Quy ước đặt tên: `T`, `U`, `V` (hoặc tên nói lên ý: `K`, `V` cho map, `E` cho element)

---

## 3. Generic Function

```go
func Sum[T Number](nums ...T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

```go
fmt.Println("Sum int:", Sum(1, 2, 3, 4))       // Sum int: 10
fmt.Println("Sum float:", Sum(1.5, 2.5))       // Sum float: 4

fmt.Println("Sum explicit:", Sum[int](1, 2, 3)) // chỉ định tường minh
```

Output thật:

```text
Sum int: 10
Sum float: 4
Sum explicit: 6
```

Constraint có sẵn của stdlib:

```go
func Max[T cmp.Ordered](first T, rest ...T) T {
    m := first
    for _, v := range rest {
        if v > m { m = v }
    }
    return m
}

fmt.Println("Max:", Max(5, 2, 9, 4))
fmt.Println("Max string:", Max("banana", "apple", "cherry"))
```

Output thật:

```text
Max: 9
Max string: cherry
```

> `cmp.Ordered` (Go 1.21+) = mọi kiểu so sánh được `<`, `>` — số, string.

---

## 4. Type Constraints

Constraint = **interface đặc biệt** giới hạn kiểu nào được dùng làm `T`.

### Constraint có sẵn

| Constraint | Từ | Bao gồm |
|---|---|---|
| `any` | stdlib | Mọi kiểu — không giới hạn |
| `comparable` | stdlib | Các kiểu so sánh được `==` (số, string, struct đủ field comparable…) |
| `cmp.Ordered` | `cmp` (Go 1.21+) | Có `<`, `>` (số, string) |

### Custom constraint — type union + tilde `~`

```go
type Number interface {
    ~int | ~int64 | ~float64
}
```

* `|` = **union** — "hoặc kiểu này, hoặc kiểu kia"
* `~int` = int **hoặc mọi kiểu có underlying type là int** (vd: `type Celsius int` ✅)
* Bỏ `~` → chỉ chấp nhận **chính xác** int, không nhận `Celsius`

```go
type Celsius int

var x Number = Celsius(35)    // ✅ OK nhờ ~int
```

---

## 5. Generic Type (struct) + Method

```go
type Stack[T any] struct {
    items []T
}

// Method của generic type cũng khai báo [T any] — dùng lại type parameter T của struct
func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T          // zero value của T — kiểu chưa biết vẫn "sống được"
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}
```

```go
var s Stack[int]        // stack cho int — với string chỉ cần Stack[string]
s.Push(10)
s.Push(20)

v, ok := s.Pop()
fmt.Println("Pop:", v, ok)
```

Output thật:

```text
Pop: 20 true
```

> ⚠️ Method **không** tự khai báo type parameter mới — chỉ dùng lại `T` của type.

---

## 6. Type Inference — compiler tự suy luận

```go
Sum(1, 2, 3, 4)        // tham số là int → compiler tự suy luận T = int
Sum[int](1, 2, 3)      // hoặc chỉ định tường minh khi cần
```

* Hầu hết các trường hợp **không cần** ghi `[int]` — compiler suy luận từ tham số
* Dùng tường minh khi: không có tham số nào mang kiểu `T`, hoặc muốn cưỡng chế kiểu khác

---

## Summary

* Generics = **1 hàm/kiểu cho nhiều kiểu dữ liệu**, kiểm tra kiểu lúc **compile** (khác `interface{}` — panic lúc runtime)
* Cú pháp: `func F[T constraint](...)`, `type S[T constraint] struct`
* **Constraint** = interface đặc biệt: `any` (mọi kiểu), `comparable` (`==`), `cmp.Ordered` (`<`, `>`), custom với union `|` + tilde `~`
* `~T` chấp cả kiểu có **underlying type** là `T`; bỏ `~` thì chỉ nhận chính xác `T`
* Generic struct có method dùng lại `[T]` của type — method không tự tạo type parameter mới
* **Type inference**: compiler tự suy luận `T` từ tham số — hiếm khi cần ghi tường minh
