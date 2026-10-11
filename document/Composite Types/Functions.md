# Functions trong Go (Golang)

7 phần dưới đây khớp 1:1 với node **Functions** của roadmap.sh.
Code mẫu: `fundamental/2.composite_type/loop_function.go`

---

## 1. Function basic

```go
func test(x, y int) int {
	return x + y
}

fmt.Println(test(2, 4))
```

* Cú pháp: `func tên(thamSo kieu) kieuTraVe { ... }` — gộp tham số: `x, y int` ≡ `x int, y int`
* Go **không** có default argument và **không** overloading (trùng tên = lỗi compile)
* Không trả gì thì để trống: `func f() {}`

---

## 2. Variadic Function

```go
func variadic(numbers ...int) {
	fmt.Println(numbers)        // bên trong, numbers là []int
}

variadic(1, 34, 54, 5, 65)     // → [1 34 54 5 65]
```

* `...int` = nhận 0 hoặc nhiều giá trị; dùng được `len(numbers)`, `append(numbers, ...)`
* Chỉ được đứng ở **tham số cuối**
* Truyền slice vào: `variadic(s...)` → "mở" slice ra thành từng phần tử

---

## 3. Multiple return value

```go
func an() (bool, int) {
	return true, 2
}

ok, giaTri := an()
```

* Trả nhiều giá trị, khai báo kiểu theo cặp; bỏ qua giá trị nào thì dùng `_`: `_, giaTri := an()`
* Pattern chuẩn của Go: `v, err := doSomething()` → luôn kiểm tra `err` trước khi dùng `v`

---

## 4. Anonymous function

```go
add := func(a, b int) int {
	return a + b
}
fmt.Println(add(3, 4))          // 7
```

* Hàm **không tên**, gán vào biến rồi gọi như hàm bình thường
* Gọi ngay tại chỗ: `(func(a, b int) int { return a + b })(3, 4)`
* Rất hay dùng làm tham số cho hàm khác (callback, `http.HandleFunc`, `sort.Slice`)

---

## 5. Closure

```go
func clouse() func() int {
	x := 10
	return func() int {
		x++        // hàm con "nhớ" biến x của hàm ngoài
		return x
	}
}

counter := clouse()
fmt.Println(counter())          // 11
fmt.Println(counter())          // 12
fmt.Println(counter())          // 13
```

* Closure = hàm **gom theo reference** biến của hàm bao ngoài → biến sống lâu hơn hàm cha
* Từ **Go 1.22**, biến vòng lặp đã *per-iteration*, nên lỗi "đặt closure trong vòng lặp lấy nhầm giá trị" gần như biến mất
* ⚠️ Closure giữ tham chiếu → capture slice/map lớn có thể gây rò rỉ bộ nhớ

---

## 6. Named return value

```go
func calculate(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b
	return sum, diff   // hoặc chỉ cần `return`
}
```

* Đặt tên cho giá trị trả về → gán thẳng được, và **`return` trần** trả về zero value hiện tại
  (vd `func example() (n int, ok bool, message string) { return }` → `0, false, ""`)
* ⚠️ Dùng tiết chế: đọc code khó đoán giá trị trả về, chỉ nên dùng khi cần `defer`/debug
* ⚠️ **Shadow**: viết `sum := ...` trong hàm này là tạo biến mới, **không** ghi vào named return

---

## 7. Call by value ⭐

Go **luôn copy tham số** khi gọi hàm — không có tham số `ref`/`out`.

```go
type NhanVien struct{ Luong int }

func suaNhanVien(luong int, nv NhanVien, giaTri map[string]int) {
	luong = 999        // chỉ sửa bản copy
	nv.Luong = 999     // struct cũng là kiểu giá trị
	giaTri["x"] = 999  // map là reference -> dùng chung
}

luong := 10
nv := NhanVien{Luong: 10}
gt := map[string]int{"x": 1}
suaNhanVien(luong, nv, gt)
```

Output thật:

```text
luong: 10
nv.Luong: 10
giaTri: map[x:999]
```

| Truyền vào | Sửa **giá trị** trong hàm | Biến gốc ngoài |
|---|---|---|
| `int`, `bool`, `string`, `struct` | ✅ | ❌ **không đổi** (bản copy) |
| `map`, `slice` | sửa **nội dung** | ✅ **đổi theo** (reference) |
| `*T` (pointer) | ✅ | ✅ đổi theo (bạn tự truyền `&x`) |

> **Tóm lại:** muốn hàm "sửa được" biến gọi → **trả về giá trị mới** hoặc truyền **pointer**.
> Đây là lý do map/slice được gọi là *reference type* dù Go không có pass by reference.

---

## Lỗi thường gặp

| Lỗi | Nguyên nhân | Cách sửa |
|-----|-------------|----------|
| `redeclared in this block` | Trùng tên hàm | Go không overloading |
| `too many arguments` / `not enough arguments` | Sai số lượng hoặc thứ tự tham số | Kiểm tra chữ ký hàm |
| `syntax error: unexpected ...` | `...` đứng không phải cuối | Đưa nó xuống tham số cuối |
| Hàm gọi xong biến gốc vẫn vậy | Go copy tham số | Xem mục 7: pointer hoặc trả về giá trị |

