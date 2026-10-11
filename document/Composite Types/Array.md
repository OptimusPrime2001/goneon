# Array

Array là tập hợp các phần tử có cùng kiểu dữ liệu và có kích thước cố định.

## Declaration

### Khai báo và khởi tạo

```go
package main

import "fmt"

func main() {
	var numbers [5]int

	fmt.Println(numbers)
}
```

Output:

```text
[0 0 0 0 0]
```

Mọi phần tử được gán giá trị zero value của kiểu dữ liệu tương ứng.

---

### Khởi tạo với giá trị ban đầu

```go
numbers := [5]int{1, 2, 3, 4, 5}
```

---

### Compiler tự suy luận kích thước

```go
numbers := [...]int{1, 2, 3, 4, 5}
```

Go sẽ tự tính length là 5.

---

## Length

Sử dụng hàm built-in `len()`.

```go
numbers := [5]int{1, 2, 3, 4, 5}

fmt.Println(len(numbers))
```

Output:

```text
5
```

---

### Length là một phần của type

```go
var a [3]int
var b [5]int
```

Hai biến trên có kiểu dữ liệu khác nhau.

```go
a = b
```

Compile Error:

```text
cannot use b as [3]int
```

Vì:

```go
[3]int != [5]int
```

---

## Access Elements

```go
numbers := [5]int{10, 20, 30, 40, 50}

fmt.Println(numbers[0])
fmt.Println(numbers[4])
```

Output:

```text
10
50
```

---

### Update Element

```go
numbers[0] = 100
```

---

## Iterating Array

### Traditional For

```go
for i := 0; i < len(numbers); i++ {
	fmt.Println(numbers[i])
}
```

---

### Range

```go
for index, value := range numbers {
	fmt.Println(index, value)
}
```

Nếu không dùng index:

```go
for _, value := range numbers {
	fmt.Println(value)
}
```

---

## Comparison

Array có thể so sánh trực tiếp nếu các phần tử bên trong có thể so sánh được.

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}

fmt.Println(a == b)
```

Output:

```text
true
```

---

### So sánh khác nhau

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 4}

fmt.Println(a == b)
```

Output:

```text
false
```

---

## Multidimensional Array

Array có thể chứa các array khác.

### 2D Array

```go
matrix := [2][3]int{
	{1, 2, 3},
	{4, 5, 6},
}
```

Biểu diễn:

```text
[
  [1 2 3]
  [4 5 6]
]
```

---

### Truy cập phần tử

```go
fmt.Println(matrix[0][1])
```

Output:

```text
2
```

---

### Duyệt 2D Array

```go
for i := range matrix {
	for j := range matrix[i] {
		fmt.Println(matrix[i][j])
	}
}
```

---
