package main

import (
	"fmt"
)

func B1() {
	type Student struct {
		Age   int
		Name  string
		Score float32
	}
	student1 := Student{
		Age:   20,
		Name:  "Trung",
		Score: 8.5,
	}
	student2 := Student{
		Age:   21,
		Name:  "Nam",
		Score: 9.0,
	}
	student3 := Student{
		Age:   22,
		Name:  "An",
		Score: 7.5,
	}
	listStudent := []Student{student1, student2, student3}
	var averageScore float32
	var maxStudent Student
	for _, student := range listStudent {
		maxStudent = student
		if student.Score > maxStudent.Score {
			maxStudent = student
		}
		sum := averageScore + float32(student.Score)
		averageScore = sum / float32(len(listStudent))
		fmt.Printf("Name: %v, Age: %v, Score: %v\n", student.Name, student.Age, student.Score)
	}
	fmt.Printf("Average Score: %v\n", averageScore)
	fmt.Printf("Top Student: %v\n", maxStudent.Name)
} 
