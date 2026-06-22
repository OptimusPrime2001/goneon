# Struct

Struct là kiểu dữ liệu dùng để nhóm nhiều field có thể thuộc các kiểu khác nhau.

---

## Declaration

### Khai báo Struct

```go
type User struct {
	ID   int
	Name string
	Age  int
}
```

---

### Khởi tạo Struct

```go
user := User{
	ID:   1,
	Name: "Trung",
	Age:  24,
}
```

---

### Truy cập field

```go
fmt.Println(user.Name)
```

Output:

```text
Trung
```

---

### Cập nhật field

```go
user.Age = 25
```

---

## Anonymous Struct

Struct không cần định nghĩa type riêng.

```go
user := struct {
	Name string
	Age  int
}{
	Name: "Trung",
	Age:  24,
}
```

Thường dùng cho:

* Test data
* Response tạm thời
* JSON mapping đơn giản

---

## Nested Struct

Struct có thể chứa struct khác.

```go
type Address struct {
	City    string
	Country string
}

type User struct {
	Name    string
	Address Address
}
```

---

### Sử dụng

```go
user := User{
	Name: "Trung",
	Address: Address{
		City:    "Hanoi",
		Country: "Vietnam",
	},
}
```

Truy cập:

```go
fmt.Println(user.Address.City)
```

Output:

```text
Hanoi
```

---

### Anonymous Nested Struct

```go
type User struct {
	Name string

	Address struct {
		City string
	}
}
```

---

## Struct Tag

Struct tag là metadata gắn vào field.

```go
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

Struct tag thường được sử dụng bởi:

* encoding/json
* validation libraries
* database ORM
* YAML libraries

---

### Ví dụ JSON

```go
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

```go
user := User{
	ID:   1,
	Name: "Trung",
}
```

Marshal:

```go
data, _ := json.Marshal(user)

fmt.Println(string(data))
```

Output:

```json
{
  "id": 1,
  "name": "Trung"
}
```

---

### Bỏ qua field

```go
type User struct {
	Password string `json:"-"`
}
```

Field sẽ không xuất hiện khi marshal.

---

### Omitempty

```go
type User struct {
	Name string `json:"name,omitempty"`
}
```

Nếu Name là zero value thì field sẽ bị bỏ qua.

---

## Struct Comparison

Struct có thể so sánh bằng `==` nếu tất cả field bên trong đều comparable.

```go
type User struct {
	ID   int
	Name string
}
```

```go
u1 := User{
	ID:   1,
	Name: "Trung",
}

u2 := User{
	ID:   1,
	Name: "Trung",
}

fmt.Println(u1 == u2)
```

Output:

```text
true
```

---

### Struct không comparable

```go
type User struct {
	Name string
	Tags []string
}
```

```go
u1 == u2
```

Compile Error:

```text
invalid operation: struct containing []string cannot be compared
```

Nguyên nhân:

* Slice không comparable
* Map không comparable
* Function không comparable

---

## Anonymous Field (Embedding)

Go hỗ trợ embedding để tạo composition.

```go
type Address struct {
	City string
}

type User struct {
	Name string
	Address
}
```

Khởi tạo:

```go
user := User{
	Name: "Trung",
	Address: Address{
		City: "Hanoi",
	},
}
```

Có thể truy cập trực tiếp:

```go
fmt.Println(user.City)
```

Thay vì:

```go
fmt.Println(user.Address.City)
```

---

# Summary

## Array

* Fixed size
* Elements cùng kiểu dữ liệu
* Comparable nếu các phần tử comparable
* Length là một phần của type
* Hỗ trợ multidimensional array

## Struct

* Nhóm nhiều field khác kiểu
* Có named struct và anonymous struct
* Hỗ trợ nested struct
* Hỗ trợ embedding
* Hỗ trợ struct tag
* Comparable nếu tất cả field đều comparable
