# Pointer trong Go (Golang)

Pointer là biến **giữ địa chỉ** của biến khác, thay vì giữ giá trị trực tiếp.
Dùng `&` để lấy địa chỉ, `*` để đọc/ghi giá trị tại địa chỉ đó.

4 phần dưới đây khớp 1:1 với node **Pointer** của roadmap.sh.
Code mẫu: `fundamental/2.composite_type/pointer.go`

---

## 1. Pointer basic

```go
i, j := 42, 1701

p := &i             // & = address operator: p giữ ĐỊA CHỈ của i
fmt.Println(p)      // 0xc000014088 (địa chỉ, thay đổi mỗi lần chạy)
```

### Dereference

```go
*p = 21             // * = dereference: ghi 21 vào giá trị mà p trỏ tới
fmt.Println(i)      // 21

p = &j              // đổi "target": p giờ trỏ sang j
*p = *p / 37        // đọc *p (=1701), chia nguyên 37, ghi lại → j = 45
fmt.Println(j)      // 45
```

* `&i` → địa chỉ của `i`, kiểu `*int` (con trỏ tới `int`)
* `*p` → giá trị **tại** địa chỉ đó; dùng được cả để đọc (`fmt.Println(*p)`) lẫn ghi (`*p = 21`)
* Đổi `p = &j` thì mọi thao tác `*p` sau đó tác động tới `j`

---

### Zero value của pointer là `nil`

```go
var p *int          // mặc định = nil, chưa trỏ tới đâu

fmt.Println(p)      // <nil>

fmt.Println(*p)     // ❌ panic
```

```text
panic: runtime error: invalid memory address or nil pointer dereference
```

Cách an toàn — luôn check trước khi deref:

```go
if p != nil {
    fmt.Println(*p)
}
```

---

### Go không có pointer arithmetic

```go
p + 1               // ❌ compile error — khác C/C++
```

Go tự quản lý bộ nhớ, không cho "trừ/cong" con trỏ. Muốn đi sang phần tử kế thì dùng slice.

---

## 2. Pointer với Struct

### Trỏ tới struct và sửa field

```go
type AppUser struct {
    fname string
    lname string
    age   int
}

appUser := AppUser{fname: "Hasin", lname: "Shekih", age: 35}

p := &appUser
p.age = 15                      // KHÔNG cần viết (*p).age

fmt.Println(appUser)            // {Hasin Shekih 15}
```

### Auto-deref khi truy cập field

Go tự giải phóng dấu `*` khi dùng `.` trên con trỏ:

```go
p.age = 15          // tự hiểu là (*p).age = 15
```

Viết `(*p).age` cũng đúng nhưng hiếm khi cần.

---

### Khi nào nên truyền struct qua pointer?

| Trường hợp | Lý do |
|---|---|
| Hàm/method cần **sửa** struct gốc | Truyền value chỉ sửa trên bản copy |
| Struct **lớn** (nhiều field) | Tránh copy toàn bộ dữ liệu mỗi lần gọi |
| Struct chứa slice/map | Muốn mọi nơi dùng chung dữ liệu |

```go
func changeAge(u *AppUser) {
    u.age = 20      // sửa trên bản gốc
}
```

> Go luôn **pass by value** — tham số struct được copy. Muốn hàm sửa được gốc thì phải truyền `&`.

---

## 3. Pointer với Slice & Map

Slice và Map là **reference type**: biến chỉ giữ *tham chiếu* tới dữ liệu gốc,
nên truyền qua hàm vẫn sửa được phần tử — **không cần** truyền `&`.

### Sửa phần tử → thấy ngay trên bản gốc

```go
person := map[string]int{"age": 25}

change := func(m map[string]int) {
    m["age"] = 200          // sửa phần tử → bản gốc bị ảnh hưởng
}
change(person)

fmt.Println(person["age"])  // 200
```

```go
numbers := []int{1, 2, 3, 4, 5}

changeSlice := func(s []int) {
    s[1] = 200              // sửa phần tử → bản gốc bị ảnh hưởng
}
changeSlice(numbers)

fmt.Println(numbers)        // [1 200 3 4 5]
```

---

### Nhưng gán lại biến local thì KHÔNG ảnh hưởng bản gốc

```go
change := func(m map[string]int) {
    m["age"] = 200                  // ✅ sửa phần tử → person bị đổi
    m = map[string]int{"age": 31}   // ❌ chỉ đổi con trỏ local, person giữ nguyên
}
```

Tóm tắt:

| Thao tác trong hàm | Ảnh hưởng bản gốc? |
|---|---|
| Sửa phần tử (`m["k"] = v`, `s[0] = v`) | ✅ Có |
| Gán lại biến (`m = ...`, `s = ...`) | ❌ Không — chỉ đổi biến local |

---

### Trường hợp đặc biệt: `append`

```go
func add(s []int) {
    s = append(s, 4)        // gán lại biến local → không return ra ngoài
}

numbers := []int{1, 2, 3}
add(numbers)

fmt.Println(numbers)        // [1 2 3]  ← VẪN LÀ [1 2 3]
```

Nguyên nhân:

* `numbers := []int{1, 2, 3}` có `len = 3`, `cap = 3`
* `append` vượt `cap` → Go tạo array mới, copy dữ liệu, gán vào biến local `s`
* Hàm không trả slice ra ngoài nên bên ngoài **không thấy** thay đổi

Muốn hàm append có hiệu lực, có 2 cách:

```go
// Cách 1: trả slice ra ngoài
func add(s []int) []int {
    return append(s, 4)
}
numbers = add(numbers)

// Cách 2: truyền *[]int
func addPtr(s *[]int) {
    *s = append(*s, 4)
}
addPtr(&numbers)
```

> ⚠️ Lưu ý: nếu `cap` còn dư thì phần tử mới được ghi vào underlying array,
> nhưng `len(numbers)` bên ngoài **vẫn không đổi** → rất dễ gây bug khó tìm.
> "Dòng máu" của Go: muốn append có hiệu lực thì luôn `s = append(s, ...)` và **gán lại / return**.

---

## 4. Pass by value vs pointer

Go **luôn** pass by value — tham số được copy trước khi vào hàm. Sự khác biệt nằm ở *cái gì bị copy*:

| | Truyền value (`T`) | Truyền pointer (`*T`) |
|---|---|---|
| Cái được copy | Toàn bộ dữ liệu | Chỉ địa chỉ (8 byte) |
| Sửa được bản gốc? | ❌ Không | ✅ Có |
| Chi phí với dữ liệu lớn | Cao | Thấp |
| Nil pointer | Không có | Deref nil → panic |

```go
func byValue(n int)    { n = 100 }      // sửa trên bản copy
func byPointer(n *int) { *n = 100 }     // sửa trên bản gốc

x := 1
byValue(x)      fmt.Println(x)      // 1   — không đổi
byPointer(&x)   fmt.Println(x)      // 100 — đã đổi
```

Gợi ý nhanh:

* Struct nhỏ, chỉ đọc → truyền value
* Cần sửa gốc, struct lớn, hoặc receiver có method mutate → truyền pointer
* Slice/map đã là reference → không cần `&` trừ khi muốn **gán lại** biến (xem mục 3)

> Pointer receiver (`func (u *AppUser) SetAge(...)`) là cách phổ biến nhất dùng pointer trong Go —
> chi tiết ở doc **Methods**.

---

## 5. Summary

* `&x` lấy địa chỉ, `*p` đọc/ghi giá trị tại địa chỉ đó — kiểu con trỏ là `*T`
* Zero value của con trỏ là `nil` — deref nil = panic, luôn check `if p != nil`
* Go **không có** pointer arithmetic (khác C/C++)
* Go luôn pass by value; struct được copy toàn bộ → muốn sửa gốc thì truyền `&`
* Auto-deref: `p.field` ≡ `(*p).field`
* Slice/Map là reference type: sửa phần tử thấy ngay, nhưng gán lại biến local không ảnh hưởng bản gốc
* `append` gán vào biến local trong hàm sẽ "mất" — phải return hoặc truyền `*[]int`
* `new(T)` trả `*T` trỏ tới zero value; `&T{...}` trả `*T` trỏ tới struct đã khởi tạo — dùng `&` là phổ biến hơn
