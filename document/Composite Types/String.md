1. Byte là gì?

- Byte là một đơn vị dữ liệu gồm 8 bit (0 và 1), có thể lưu giá trị từ `0` đến `255: 2^8`
- Là bí danh của UNIT8

Trong bảng mã ACII, ta có thể tra ra 1 kí tự tương ứng với DEC, HEX là bao nhiêu
VD: A → 65 ( DEC ) → 41 ( HEX ) 

A = 1 byte = 0100 0001

$$
\begin{aligned} \text{Dec} &= (0 \times 2^7) + (1 \times 2^6) + (0 \times 2^5) + (0 \times 2^4) + (0 \times 2^3) + (0 \times 2^2) + (0 \times 2^1) + (1 \times 2^0) \\ &= 0 + 64 + 0 + 0 + 0 + 0 + 0 + 1 \\ &= \mathbf{65} \end{aligned}
$$

**Đổi sang HEX :** 

**Chia nhóm 4 bit:** `0100` và `0001`
• Nhóm thứ nhất `0100` = $0 \times 8 + 1 \times 4 + 0 \times 2 + 0 \times 1 = \mathbf{4}$
• Nhóm thứ hai `0001` = $0 \times 8 + 0 \times 4 + 0 \times 2 + 1 \times 1 = \mathbf{1}$ 

Kêt quả là **`41`** (thường viết là `0x41`).

**Ví dụ trong GO**

```
var b dec = 65
fmt.Println(b)         // 65
fmt.Printf("%c\n", b)  // A
```

Tại sao lại in ra `A`?

- `b` đang lưu số `65`.
- Theo bảng mã ASCII, số `65` tương ứng với chữ `A`.
- `%c` yêu cầu Go hiển thị giá trị dưới dạng ký tự.

## 2. String là gì?

- String là một chuỗi byte. Trong Go, string không thể bị sửa đổi trực tiếp sau khi được tạo.
- Các kí tự có dấu thường có nhiều hơn 1 byte
- `len(s)` đếm byte, không phải lúc nào cũng đếm số ký tự.

### 3. UTF-8 là gì?

UTF-8 là một cách mã hóa ký tự Unicode thành các byte.

Máy tính lưu và xử lý dữ liệu dưới dạng số. Vì vậy, để lưu các ký tự như `A`, `ế`, `日` hay `⌘`, chúng cần được biểu diễn thành các byte.

UTF-8 quy định cách biểu diễn đó.

| Kí tự  | byte | Hex |
| --- | --- | --- |
| A | 1 byte | Hexadecimal: `41` |
| ế | 3 byte trong UTF-8 | Hexadecimal: `e1 ba bf` |
| 日 | 3 byte trong UTF-8 | Hexadecimal: `e6 97 a5` |

```go
fmt.Println(len("A"))  // 1
fmt.Println(len("ế"))  // 3
fmt.Println(len("日")) // 3
```

### 4. Unicode và code point là gì?

Để hiểu UTF-8, bạn cần biết thêm Unicode.

- Unicode là tiêu chuẩn gán mã cho các ký tự và nhiều thành phần văn bản.
- Code point là một giá trị cụ thể trong Unicode, thường được viết dưới dạng `U+` kèm mã hexadecimal.
- Ví dụ: `A` có code point `U+0041`; `日` có code point `U+65E5`.

Điểm khác biệt:

- Unicode xác định mã của ký tự.
- UTF-8 quy định cách mã đó được biểu diễn thành byte.

### 5. Rune là gì?

Trong Go, `rune` là bí danh của `int32`, dùng để biểu diễn một Unicode code point.

```
var r rune = '日'

fmt.Println(r)         // 26085
fmt.Printf("%c\n", r)  // 日
fmt.Printf("%U\n", r)  // U+65E5
```

Tại sao `fmt.Println(r)` in ra `26085`?

Vì `r` lưu giá trị số của code point `U+65E5`, tương đương số thập phân `26085`.

Các cách in khác nhau:

- `%d`: in giá trị số thập phân.
- `%c`: hiển thị ký tự tương ứng.
- `%U`: hiển thị mã Unicode.

## Byte khác rune như thế nào?

Lấy ký tự `日` làm ví dụ:

```
s := "日"

fmt.Println(len(s)) // 3

for _, r := range s {
    fmt.Println(r) // 26085
}
```

- `len(s) == 3` vì UTF-8 dùng 3 byte để mã hóa `日`.
- Vòng `range` chạy một lần vì chuỗi có một rune.
- Giá trị rune là `26085`, tức code point `U+65E5`.

Vậy nên, một rune có thể được mã hóa thành nhiều byte.

Ghi nhớ:

- `for` với `s[i]`: duyệt từng byte.
- `for range s`: duyệt từng rune.
- `index` trong `range` không phải số thứ tự ký tự.

### Các format thường dùng trong `fmt`

Giả sử:

```
s:="Hi 日"
```

| Format | Ví dụ | Ý nghĩa |
| --- | --- | --- |
| `%s` | `fmt.Printf("%s", s)` | In string bình thường |
| `%q` | `fmt.Printf("%q", s)` | In string trong dấu ngoặc kép, escape ký tự đặc biệt |
| `%x` | `fmt.Printf("%x", s)` | In các byte dưới dạng hex, viết liền |
| `% X` | `fmt.Printf("% X", s)` | In các byte hex, cách nhau bằng dấu cách |
| `%+q` | `fmt.Printf("%+q", s)` | In string có dấu ngoặc kép, escape ký tự Unicode không phải ASCII |
| `%v` | `fmt.Printf("%v", s)` | In giá trị theo format mặc định |
| `%#v` | `fmt.Printf("%#v", s)` | In biểu diễn giá trị theo cú pháp Go |

### Libraries

Package `fmt` chỉ giúp hiển thị dữ liệu. Nếu muốn xác thực UTF-8, hãy dùng `unicode/utf8`:

```
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    s := "Hello 日"

    fmt.Println("String:", s)
    fmt.Println("Bytes:", len(s))
    fmt.Println("Runes:", utf8.RuneCountInString(s))
    fmt.Println("Valid UTF-8:", utf8.ValidString(s))
}
```

Kết quả:

```
String: Hello 日
Bytes: 9
Runes: 7
Valid UTF-8: true
```

Vì `"Hello 日"` có 6 byte cho `Hello`  và 3 byte cho `日`, tổng cộng 9 byte; nhưng chỉ có 7 rune.