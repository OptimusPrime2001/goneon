Dưới đây là nội dung tổng hợp đầy đủ về **Slice trong Go** để bạn có thể copy vào file `Slice.md` (bao gồm cơ bản → nâng cao → phỏng vấn → tips thực tế).

---

# 📘 Slice trong Go (Golang)

## 1. Khái niệm cơ bản

Slice là một **cấu trúc dữ liệu động (dynamic array)** trong Go, được xây dựng trên mảng (array).

```go
var s []int
```

Slice không lưu dữ liệu trực tiếp mà **tham chiếu tới một array bên dưới (underlying array)**.

---

## 2. Cấu trúc của Slice

Một slice gồm 3 thành phần:

* Pointer → trỏ đến array
* Length → số phần tử hiện tại
* Capacity → sức chứa tối đa từ vị trí bắt đầu

```go
type slice struct {
    ptr *T
    len int
    cap int
}
```

---

## 3. Khởi tạo Slice

### Cách 1: khai báo trực tiếp

```go
s := []int{1, 2, 3}
```

### Cách 2: dùng `make`

```go
s := make([]int, 5)      // len = 5, cap = 5
s := make([]int, 3, 10)  // len = 3, cap = 10
```

### Cách 3: từ array

```go
arr := [5]int{1, 2, 3, 4, 5}
s := arr[1:4]
```

---

## 4. Len và Cap

```go
s := []int{1, 2, 3}

fmt.Println(len(s)) // 3
fmt.Println(cap(s)) // phụ thuộc underlying array
```

---

## 5. Cách hoạt động của Slice (quan trọng)

Slice **không copy dữ liệu**, chỉ trỏ tới array gốc.

```go
arr := []int{1, 2, 3, 4}
s := arr[1:3]

s[0] = 100

fmt.Println(arr) // [1 100 3 4]
```

👉 Thay đổi slice → ảnh hưởng array gốc

---

## 6. Append trong Slice

```go
s := []int{1, 2, 3}
s = append(s, 4)
```

### Khi vượt capacity

Go sẽ:

* tạo array mới
* copy dữ liệu
* tăng capacity (thường x2)

---

## 7. Copy Slice

```go
src := []int{1, 2, 3}
dst := make([]int, len(src))

copy(dst, src)
```

👉 `copy` giúp tránh dùng chung underlying array

---

## 8. Nil Slice vs Empty Slice

### Nil slice

```go
var s []int
```

* s == nil → true
* len = 0, cap = 0

### Empty slice

```go
s := []int{}
```

* s != nil
* len = 0

👉 Phỏng vấn hay hỏi điểm này

---

## 9. Slice operations (cắt slice)

```go
s := []int{1, 2, 3, 4, 5}

fmt.Println(s[1:4]) // [2 3 4]
fmt.Println(s[:3])  // [1 2 3]
fmt.Println(s[2:])  // [3 4 5]
```

---

## 10. Delete phần tử trong Slice

Go không có delete built-in:

```go
s := []int{1, 2, 3, 4}

// xóa index 1
s = append(s[:1], s[2:]...)
```

---

## 11. Insert vào Slice

```go
s := []int{1, 3, 4}
s = append(s[:1], append([]int{2}, s[1:]...)...)
```

---

## 12. Slice tricks (quan trọng thực tế)

### 12.1 Tránh memory leak khi slice lớn

```go
small := big[100:200]
```

👉 vẫn giữ reference toàn bộ `big`

Cách fix:

```go
small := append([]int(nil), big[100:200]...)
```

---

### 12.2 Pre-allocate để tối ưu performance

```go
s := make([]int, 0, 1000)
```

👉 tránh realloc nhiều lần khi append

---

### 12.3 Check capacity growth

```go
for i := 0; i < 10; i++ {
	s = append(s, i)
	fmt.Println(len(s), cap(s))
}
```

---

## 13. Slice trong function (quan trọng phỏng vấn)

Slice là **reference type**, nhưng header được copy.

```go
func update(s []int) {
	s[0] = 100
}
```

👉 thay đổi vẫn ảnh hưởng bên ngoài

Nhưng:

```go
func appendTest(s []int) {
	s = append(s, 10)
}
```

👉 không ảnh hưởng bên ngoài (nếu không return)

---

## 14. Pass by value vs reference (câu hỏi phỏng vấn)

Slice:

* header → pass by value
* underlying array → shared

👉 Đây là câu hỏi rất hay gặp

---

## 15. So sánh Slice

❌ Không thể so sánh slice bằng `==` (trừ nil)

```go
// sai
s1 == s2
```

✔ đúng cách:

```go
import "reflect"

reflect.DeepEqual(s1, s2)
```

---

## 16. Range với Slice

```go
for i, v := range s {
	fmt.Println(i, v)
}
```

---

## 17. Slice vs Array

| Array      | Slice          |
| ---------- | -------------- |
| Fixed size | Dynamic        |
| Value type | Reference-like |
| Hiếm dùng  | Rất phổ biến   |

---

## 18. Câu hỏi phỏng vấn thường gặp

### 1. Slice là gì?

→ struct gồm ptr, len, cap trỏ tới array

---

### 2. Slice có phải reference type không?

→ Không hoàn toàn, là value type nhưng chứa pointer

---

### 3. append hoạt động như thế nào?

→ Nếu đủ cap → dùng array cũ
→ Nếu không đủ → tạo array mới

---

### 4. Vì sao append trong function không thay đổi ngoài?

→ vì slice header bị copy

---

### 5. Difference giữa nil slice và empty slice?

→ nil slice == nil, empty slice != nil

---

## 19. Tips thực tế (rất quan trọng)

* Luôn dùng `make` khi biết trước size
* Tránh giữ reference slice từ data lớn
* Luôn return slice nếu dùng append trong function
* Cẩn thận với shared underlying array
* Dùng `copy` khi cần isolation data

---

## 20. Summary nhanh

Slice trong Go là:

* Dynamic array
* Có len + cap
* Dùng underlying array
* append có thể realloc
* reference behavior dễ gây bug nếu không hiểu
