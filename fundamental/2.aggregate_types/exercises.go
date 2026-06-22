// Student Management
// Employee Information
// Nested Address Struct
package main

import (
	"encoding/json"
	"fmt"
)

func StudenManagement() {
	type Student struct {
		Id   int
		Name string
		Age  int
		GPA  float64
	}
	listStudent := [2]Student{
		{
			Id: 20, Name: "Lê Văn Trung", Age: 20, GPA: 3.3,
		},
		{Id: 12, Name: "Nguyễn Việt Hung", Age: 20, GPA: 3.5},
	}
	fmt.Printf("%-5s %-25s %-5s %-5s\n", "ID", "Name", "Age", "GPA")
	if listStudent[0] == listStudent[1] {
		fmt.Println("Hai học sinh có cấu hình giống nhau")
	} else {
		fmt.Println("Hai học sinh có cấu hình khác nhau")
	}
	maxGpaStudent := listStudent[0]
	for _, std := range listStudent {
		if std.GPA > maxGpaStudent.GPA {
			maxGpaStudent = std
		}
		std.GPA += 0.2
		fmt.Printf(
			"%-5d %-25s %-5d %-5.2f\n",
			std.Id,
			std.Name,
			std.Age,
			std.GPA,
		)
	}
	fmt.Println("Top Student", maxGpaStudent.Name)
}

func EmployeeInformation() {
	employee := struct {
		Id         int     `json:"id"`
		Name       string  `json:"name"`
		Department string  `json:"department"`
		Salary     float64 `json:"salary"`
	}{
		Id:         1,
		Name:       "Pedro",
		Department: "HAUI",
		Salary:     30_000_000,
	}
	jsonEmployee, err := json.Marshal(employee)
	if err == nil {
		fmt.Printf("Thông tin employee: %s\n", string(jsonEmployee))
	}
}
