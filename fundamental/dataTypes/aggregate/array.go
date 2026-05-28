package aggregate

import (
	"fmt"
)

func Array() {
	// 1. Nhóm array ( Array) : Mảng có kích thước cố định và lưu các phần tử cùng cùng kiểu dữ liệu
	var a [3]int               // [0,0,0]
	b := [7]int{1, 2, 3, 4, 5} // [1,2,3,4,5,0,0]
	c := [...]string{"a", "b", "c"}
	fmt.Println("Array values ", a, b, c[2])
	fmt.Println("Array length ", len(a), len(b), len(c))
}