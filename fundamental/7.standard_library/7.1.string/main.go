package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func package_strings() {
	s := "Hello, Go world"

	// Kiểm tra và tìm kiếm

	strings.Contains(s, "Go")                // true
	strings.ContainsAny(s, "xyz")            // false (chứa bất kỳ ký tự nào trong tập)
	strings.ContainsRune(s, 'G')             // true
	strings.HasPrefix(s, "Hello")            // true
	strings.HasSuffix(s, "World")            // true
	strings.Index(s, "Go")                   // 7 (vị trí byte, -1 nếu không có)
	strings.LastIndex(s, "o")                // 13
	strings.IndexByte(s, ',')                // 5
	strings.Count(s, "o")                    // 3
	strings.EqualFold("Go", "GO")            // true (so sánh không phân biệt hoa thường)
	strings.ContainsFunc(s, unicode.IsDigit) // false (Go 1.21+)

	// Tách và nối
	strings.Split("a,b,c", ",")                  // ["a" "b" "c"]
	strings.SplitN("a,b,c,d", ",", 2)            // ["a" "b,c,d"]
	strings.Join([]string{"a", "b", "c"}, "-")   // "a-b-c"
	strings.Fields("  a   b \t c\n")             // ["a" "b" "c"] (tách theo mọi khoảng trắng)
	strings.FieldsFunc("a1b2c", unicode.IsDigit) // ["a" "b" "c"]

	//Cut: cách tách hai nửa gọn nhất
	strings.Cut("key=value", "=")
	// before="key", after="value", found=true

	strings.CutPrefix("Bearer abc123", "Bearer ") // rest="abc123", ok=true
	strings.CutSuffix("main.go", ".go")           // name="main", ok=true (Go 1.20+)

	//Cắt khoảng trắng và ký tự
	strings.TrimSpace("  hi \n")                   // "hi"
	strings.Trim("xxhixx", "x")                    // "hi"
	strings.TrimLeft("xxhixx", "x")                // "hixx"
	strings.TrimRight("xxhixx", "x")               // "xxhi"
	strings.TrimPrefix("prefix-body", "prefix-")   // "body"
	strings.TrimSuffix("file.txt", ".txt")         // "file"
	strings.TrimFunc("123abc456", unicode.IsDigit) // "abc"

	// Biến đổi
	strings.ToUpper("go")                 // "GO"
	strings.ToLower("GO")                 // "go"
	strings.Repeat("ab", 3)               // "ababab"
	strings.Replace("a.b.c", ".", "-", 1) // "a-b.c" (thay tối đa n lần)
	strings.ReplaceAll("a.b.c", ".", "-") // "a-b-c"
	strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		} // trả -1 để xóa ký tự
		return unicode.ToUpper(r)
	}, "a b c") // "ABC"

	// Hành vi đặc biệt
	strings.Split("", ",")   // [""]  → độ dài 1, KHÔNG phải rỗng
	strings.Fields("")       // []    → độ dài 0
	strings.Split("abc", "") // ["a" "b" "c"] (tách theo từng ký tự UTF-8)
	strings.Index("abc", "") // 0

	// strings.Builder: nối chuỗi hiệu quả
	var sb strings.Builder
	sb.Grow(64) // tùy chọn: cấp phát trước nếu ước lượng được kích thước
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "item-%d;", i)
	}
	sb.WriteString("end")
	sb.WriteByte('!')
	sb.WriteRune('✓')
	fmt.Println(sb.String(), sb.Len())
}

func package_strconv() {
	//Chuỗi → số
	strconv.Atoi("123")             // int
	strconv.ParseInt("-42", 10, 64) // (chuỗi, cơ số, số bit)
	strconv.ParseUint("42", 10, 64)
	strconv.ParseFloat("3.14", 64)
	strconv.ParseBool("true")       // chấp nhận: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
	strconv.ParseInt("ff", 16, 64)  // 255
	strconv.ParseInt("0x1f", 0, 64) // 31
	strconv.ParseInt("300", 10, 8)  // lỗi: out of range (int8 tối đa 127)

	// Số → chuỗi
	strconv.Itoa(123)                         // "123"
	strconv.FormatInt(255, 2)                 // "11111111" (nhị phân)
	strconv.FormatInt(255, 16)                // "ff"
	strconv.FormatFloat(3.14159, 'f', 2, 64)  // "3.14"
	strconv.FormatFloat(3.14159, 'f', -1, 64) // "3.14159" (-1 = ít chữ số nhất để giữ chính xác)
	strconv.FormatBool(true)                  // "true"
}
func main() {
	s := "Xin chào"
	fmt.Println(len(s))                    // 9 (vì "à" chiếm 2 byte)
	fmt.Println(utf8.RuneCountInString(s)) // 8 ký tự
	fmt.Println(len([]rune(s)))            // 8 (cách khác, nhưng cấp phát thêm)
	for id, r := range s {
		fmt.Printf("%d: %c \n", id, r) // i là vị trí byte, r là rune
	}

	r := []rune("Xin chào")
	fmt.Println(string(r[:3])) // "Xin"
	package_strings()
}
