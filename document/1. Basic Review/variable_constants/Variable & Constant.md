### `var`

`var` dùng để khai báo biến theo cách minh tường. Bạn có thể khai báo kèm kiểu dữ liệu hoặc để Go tự suy luận giá trị ban đầu.

```go
package main

import "fmt"

func main() {
    var age int = 25
    var name = "Go"

    fmt.Println(age)
    fmt.Println(name)
}
```

Khi nên dùng `var`:

- Khi muốn khai báo biến ở `package scope`.
- Khi muốn dùng giá trị mặc định của biến.
- Khi cần viết rõ ràng kiểu dữ liệu.

### `:=`

`:=` là cách khai báo biến ngắn gọn, chỉ dùng được bên trong `function`. Go sẽ tự động suy luận kiểu từ giá trị bên phải.

```go
package main

import "fmt"

func main() {
    count := 10
    title := "basic review"

    fmt.Println(count)
    fmt.Println(title)
}
```

Lưu ý:

- `:=` không dùng được ngoài `function`.
- Ít nhất phải có 1 biến mới ở vế trái.

### `const`

`const` dùng để khai báo hằng số, nghĩa là giá trị không thay đổi sau khi được gán.

```go
package main

import "fmt"

const Pi = 3.14159
const AppName = "goneon"

func main() {
	const (
		Pending = iota // 0
		Processing      // 1
		Done            // 2
	)
	fmt.Println(Pi)
	fmt.Println(AppName)
}
```

Nên dùng `const` cho:

- Các giá trị cố định như PI, `port` mặc định, `status code` tự định nghĩa.
- Các chuỗi cấu hình không thay đổi trong chương trình.

`iota` trong Go để :

- Dùng trong const để tự động tăng giá trị
- Bắt đầu từ 0, mỗi dòng tăng lên 1
- Thường dùng để tạo enum