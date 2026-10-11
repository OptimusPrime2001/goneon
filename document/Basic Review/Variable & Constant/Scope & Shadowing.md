# Scope & Shadowing

Mục tiêu của phần này là hiểu rõ biến trong Go được truy cập từ đâu và chuyện gì xảy ra khi khai báo biến trùng tên ở các tầng lồng nhau. Sau khi học xong, bạn cần nắm được:

- Scope là gì và Go có những loại scope nào.
- Phân biệt `universe scope`, `package scope`, `file scope`, `function scope` và `block scope`.
- Shadowing hoạt động như thế nào.
- Cạm bẫy thường gặp với `err` bị shadow và cách tránh.

***

## 1. Scope là gì?

`Scope` (vùng scope) là vùng code mà một biến có thể được truy cập. Biến chỉ tồn tại và hữu dụng trong phạm vi nó được khai báo. Ra khỏi vùng đó, biến không còn tồn tại nữa.

```go
package main

import "fmt"

func main() {
	x := 10 // x chỉ tồn tại trong main
	fmt.Println(x)
}

// fmt.Println(x) // ❌ lỗi: x không còn tồn tại ở đây
```

Ý nghĩa:

- Giúp quản lý vòng đời của biến.
- Tránh xung đột tên giữa các vùng code khác nhau.
- Giúp compiler bắt lỗi khi bạn truy cập biến "ma".

***

## 2. Các loại Scope trong Go

Go có 5 loại scope, xếp theo thứ tự từ lớn đến nhỏ:

| Thứ tự | Loại scope       | Phạm vi                                    |
|--------|------------------|---------------------------------------------|
| 1      | Universe scope   | Các identifier có sẵn của Go               |
| 2      | Package scope    | Mọi file trong cùng package                |
| 3      | File scope       | `import` — chỉ trong file khai báo          |
| 4      | Function scope   | Bên trong hàm                               |
| 5      | Block scope      | Bên trong `{ }` lồng                        |

### Universe scope

`Universe scope` chứa các identifier có sẵn của Go mà không cần khai báo: kiểu dữ liệu (`int`, `string`, `bool`...), hằng số (`true`, `false`, `nil`), hàm (`len`, `cap`, `make`, `new`...).

```go
package main

import "fmt"

func main() {
	x := 10          // int đến từ universe scope
	s := "hello"     // string đến từ universe scope
	b := true        // bool đến từ universe scope

	fmt.Println(x, s, b)
}
```

Đặc điểm:

- Luôn tồn tại, không cần `import` hay khai báo.
- Có thể bị shadow: khai báo `int` hay `len` trong hàm sẽ che chúng đi (nhưng **không nên**).

### Package scope

Biến khai báo **bên ngoài** mọi hàm thuộc `package scope`. Mọi file trong cùng package đều truy cập được.

```go
package main

import "fmt"

var AppName = "goneon" // package scope

func main() {
	fmt.Println(AppName) // ✅ truy cập được
}

func printName() {
	fmt.Println(AppName) // ✅ truy cập được
}
```

Đặc điểm:

- Tồn tại trong suốt vòng đời chương trình.
- Được khai báo bằng `var` hoặc `const` (không dùng `:=`).
- Name ở package scope có thể bị shadow bởi biến cùng tên trong hàm.

### File scope

`File scope` là phạm vi **chỉ trong 1 file**. Trong Go, chủ yếu là các `import` — tên package import chỉ thấy được trong file khai báo, không tự động lan sang file khác trong cùng package.

```go
// file: main.go
package main

import "fmt" // file scope — chỉ file này dùng được "fmt"

func main() {
	fmt.Println("hello")
}
```

```go
// file: util.go
package main

// fmt.Println("hi") // ❌ lỗi: "fmt" không được import ở file này
```

Đặc điểm:

- `import` ở đầu file thuộc `file scope`.
- Muốn dùng `fmt` trong `util.go` phải tự `import "fmt"` ở file đó.
- Biến khai báo bằng `var` **ở top level** vẫn là `package scope` (dù nằm trong 1 file) — nó không phải file scope.

Phân biệt nhanh:

- `import "fmt"` → **file scope** (riêng file).
- `var AppName = "..."` ở top level → **package scope** (chung mọi file).

### Function scope

Biến khai báo **bên trong hàm** thuộc `function scope`. Chỉ hàm đó truy cập được.

```go
package main

import "fmt"

func greet() {
	message := "hello" // function scope
	fmt.Println(message)
}

func main() {
	greet()
	// fmt.Println(message) // ❌ lỗi: message không tồn tại ở đây
}
```

Đặc điểm:

- Biến biến mất ngay khi hàm kết thúc.
- Dùng `:=` hoặc `var` bên trong hàm.
- Không thể truy cập từ hàm khác.

### Block scope

Mọi block `{ }` đều tạo ra một scope riêng: `if`, `for`, `switch`, `select`, hay block `{ }` tự viết.

```go
package main

import "fmt"

func main() {
	x := 100 // function scope của main

	if x > 50 {
		y := x * 2 // block scope của if
		fmt.Println(y)
	}

	// fmt.Println(y) // ❌ lỗi: y không tồn tại ngoài block if

	for i := 0; i < 3; i++ { // i chỉ tồn tại trong block for
		fmt.Println(i)
	}
	// fmt.Println(i) // ❌ lỗi
}
```

Ý nghĩa:

- Biến trong block không "leak" ra ngoài.
- Biến ngoài block vẫn đọc được bên trong (nếu không bị shadow).

### Tóm tắt phạm vi

| Loại scope      | Khai báo ở                  | Truy cập từ                                  |
|-----------------|-----------------------------|-----------------------------------------------|
| Universe scope  | Không cần khai báo          | Mọi nơi                                      |
| Package scope   | Ngoài mọi hàm (top level)   | Mọi file trong package                       |
| File scope      | `import` ở đầu file         | Chỉ file khai báo                            |
| Function scope  | Bên trong hàm               | Chỉ trong hàm đó                             |
| Block scope     | Bên trong `{ }` lồng        | Chỉ trong block đó và block con              |

***

## 3. Shadowing là gì?

`Shadowing` (che khuất) xảy ra khi bạn khai báo một biến **cùng tên** với biến ở scope bên ngoài, trong một scope lồng ở bên trong. Biến bên ngoài vẫn tồn tại nhưng bị "che" — code bên trong block sẽ dùng biến mới, không phải biến cũ.

```go
package main

import "fmt"

func main() {
	x := 10 // biến ngoài

	{
		x := 99 // biến mới, shadow biến ngoài
		fmt.Println("trong block:", x) // 99
	}

	fmt.Println("ngoài block:", x) // 10 — không đổi
}
```

Kết quả:

```text
trong block: 99
ngoài block: 10
```

Ý nghĩa:

- Biến ngoài **không bị xóa hay thay đổi**, chỉ bị che khuất tạm thời trong block.
- Ra khỏi block, biến ngoài hoạt động bình thường trở lại.
- Shadowing **không phải lỗi compile**, nên rất dễ vô tình tạo bug.

***

## 4. Các trường hợp Shadowing phổ biến

### Shadowing trong `if`

```go
package main

import "fmt"

func main() {
	x := 10

	if x > 5 {
		x := 200 // biến mới trong block if
		fmt.Println("if:", x)
	}

	fmt.Println("main:", x)
}
```

```text
if: 200
main: 10
```

### Shadowing trong `for`

```go
package main

import "fmt"

func main() {
	i := 100

	for i := 0; i < 3; i++ { // i mới, shadow i của main
		fmt.Println("for:", i)
	}

	fmt.Println("main:", i)
}
```

```text
for: 0
for: 1
for: 2
main: 100
```

### Shadowing trong `switch`

```go
package main

import "fmt"

func main() {
	x := 2

	switch x {
	case 2:
		x := 99 // shadow trong case block
		fmt.Println("case:", x)
	}

	fmt.Println("main:", x)
}
```

```text
case: 99
main: 2
```

### Cạm bẫy nổi tiếng: `err` bị shadow

Đây là bug thường gặp nhất trong Go. Khi viết `if err := ...` lồng nhau, biến `err` bên ngoài bị biến `err` bên trong shadow, khiến `err` ở ngoài không được gán.

```go
package main

import (
	"errors"
	"fmt"
)

func doWork() error {
	err := errors.New("lỗi ngoài")

	if err := innerWork(); err != nil { // ⚠️ khai báo err MỚI, shadow err ngoài
		return err
	}

	// err ở đây vẫn là "lỗi ngoài" — không phải kết quả innerWork
	fmt.Println("err ngoài:", err)

	return nil
}

func innerWork() error {
	return errors.New("lỗi trong")
}

func main() {
	fmt.Println(doWork())
}
```

```text
err ngoài: lỗi ngoài
<lỗi trong>
```

**Vì sao bug này nguy hiểm?**

- `err` bên trong `if err := ...` chỉ sống trong block `if`.
- `err` khai báo trước đó **không hề được cập nhật**.
- Code có thể bỏ qua lỗi một cách âm thầm.

### Cách sửa

Cách 1: khai báo `err` trước, tái sử dụng bằng `=` thay vì `:=`

```go
var err error
if err = innerWork(); err != nil {
	return err
}
```

Cách 2: nếu biến `err` bên ngoài đã tồn tại, dùng nhiều biến trái để tránh shadow

```go
if _, err = innerWork(); err != nil {
	return err
}
```

Cách 3: tránh khai báo `err` trùng tên ở nhiều tầng

```go
if errOuter := outerWork(); errOuter != nil {
	return errOuter
}
if errInner := innerWork(); errInner != nil {
	return errInner
}
```

***

## 5. Short variable declaration `:=` và Scope

`:=` có 2 hành vi quan trọng liên quan tới scope:

### Cùng scope — biến cũ được **reuse** (không phải shadow)

```go
package main

import "fmt"

func main() {
	x := 10
	fmt.Println(x)

	x := 20 // ❌ lỗi compile: "no new variables on left side of :="
	fmt.Println(x)
}
```

- Khi `:=` ở **cùng scope** mà không có biến mới nào bên trái → lỗi compile.
- Nếu có 1 biến mới + 1 biến cũ → biến cũ được **gán lại**, biến mới được tạo:

```go
package main

import "fmt"

func main() {
	x := 10
	x, y := 20, 30 // x được gán lại, y là biến mới

	fmt.Println(x) // 20
	fmt.Println(y) // 30
}
```

### Scope khác — biến mới được tạo (shadowing)

```go
package main

import "fmt"

func main() {
	x := 10

	if true {
		x := 99 // scope khác → tạo biến MỚI → shadow
		fmt.Println(x)
	}

	fmt.Println(x) // 10
}
```

Tóm tắt nhanh:

- `:=` cùng scope → biến cũ **reuse** (gán lại).
- `:=` scope khác → biến mới **shadow** biến cũ.

***

## 6. Lợi ích & Rủi ro

### Lợi ích của shadowing

- Tránh xung đột tên: dùng `err`, `i`, `ctx` quen thuộc mà không sợ đụng biến toàn cục.
- Code ngắn gọn hơn, không cần nghĩ tên riêng cho từng biến tạm.

### Rủi ro của shadowing

| Rủi ro                  | Mô tả                                              |
|-------------------------|----------------------------------------------------|
| Bug `err` bị nuốt       | Lỗi không được kiểm tra do `err` bên ngoài bị shadow |
| Biến ngoài không cập nhật | Ghi `x := ...` nhưng mong `x` ngoài đổi theo        |
| Khó debug               | Đọc code không rõ đang dùng biến nào               |
| Vòng lặp ngoài bị che   | `i` trong `for` che `i` của hàm                    |

***

## 7. Best Practices

- **Không đặt trùng tên** biến ở scope lồng nhau trừ khi cố ý shadow.
- **Đặt tên biến rõ ràng** cho biến mới: `errInner`, `errOuter`, `userID`, `orderID`.
- **Cảnh giác với `if err := ...`** lồng nhau — đọc kỹ biến `err` nào đang được kiểm tra.
- Dùng **`var` + `=`** khi muốn tái sử dụng biến đã khai báo ở scope khác.
- Khi code có nhiều tầng lồng, **tách hàm nhỏ** để giảm số tầng scope.
- Dùng linter (`golangci-lint` với rule `govet shadow`) để bắt shadowing tiềm ẩn.

***

## 8. Tóm tắt nhanh

- `scope` là vùng code mà biến có thể truy cập.
- Go có 5 loại scope: `universe`, `package`, `file`, `function` và `block`.
- `universe scope` gồm các kiểu và hàm có sẵn (`int`, `len`, `make`...).
- `package scope` là biến khai báo ở top level — chung mọi file trong package.
- `file scope` chủ yếu là `import` — chỉ file khai báo mới dùng được.
- `block scope` đến từ `{ }` của `if`, `for`, `switch` và block tự viết.
- `shadowing` là biến bên trong cùng tên che khuất biến bên ngoài.
- Biến ngoài **không bị đổi**, chỉ bị che tạm trong block.
- `:=` cùng scope → reuse; `:=` scope khác → shadow.
- Cạm bẫy lớn nhất: `if err := ...` lồng nhau làm `err` bên ngoài không cập nhật.
- Tránh đặt trùng tên, dùng tên biến riêng biệt cho từng tầng.

***

## 9. Exercises

### 1. Dự đoán kết quả

Không chạy code, viết ra kết quả in ra màn hình:

```go
package main

import "fmt"

func main() {
	x := 1

	if x > 0 {
		x := 10
		fmt.Println("if:", x)
	}

	for i := 0; i < 2; i++ {
		x := i
		fmt.Println("for:", x)
	}

	fmt.Println("main:", x)
}
```

Yêu cầu:

- Viết kết quả ra giấy.
- Sau đó chạy `go run` để đối chiếu.

### 2. Tìm và sửa bug

Chương trình sau in ra kết quả sai. Tìm nguyên nhân do shadowing và sửa lại:

```go
package main

import (
	"errors"
	"fmt"
)

func validate(age int) error {
	err := errors.New("age must be >= 0")

	if err := checkRange(age); err != nil {
		fmt.Println("phát hiện lỗi:", err)
	}

	// mong muốn: nếu có lỗi thì in ra, không thì in "ok"
	if err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("lỗi:", err)
	}

	return nil
}

func checkRange(age int) error {
	if age < 0 {
		return errors.New("age < 0")
	}
	return nil
}

func main() {
	fmt.Println(validate(-5))
}
```

Gợi ý:

- Xem lại biến `err` nào đang được kiểm tra ở dòng `if err == nil`.

### 3. Viết lại không dùng shadowing

Viết lại hàm sau mà **không để biến `i` bị shadow**, đồng thời giữ nguyên kết quả:

```go
package main

import "fmt"

func main() {
	i := 100

	for i := 0; i < 3; i++ {
		fmt.Println("trong for:", i)
	}

	fmt.Println("ngoài for:", i)
}
```

Mục tiêu bài này là luyện:

- Phân biệt `block scope` với `function scope`.
- Hiểu `:=` reuse khác `:=` shadow thế nào.
- Nhận ra bug `err` kinh điển trong code thật.
