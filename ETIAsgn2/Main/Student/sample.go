package student

import (
	"fmt"
)

type API interface {
	getAllStudentNames() []string
}

type FakeAPI2 struct{}

func (api FakeAPI2) getAllStudentNames() []string {
	return []string{
		"d", "e",
	}
}

type FakeAPI1 struct{}

func (fakeApi FakeAPI1) getAllStudentNames() []string {
	return []string{
		"a", "b",
	}
}

func PrintAllStudentNames(api API) {
	fmt.Printf("%v", api.getAllStudentNames())
}

// Version 1
// func GetAllStudents() string {
// 	return "a"
// }

//Version 2
func GetAllStudents() string {
	return "b"
}

func main() {
	database := FakeAPI1{}
	PrintAllStudentNames(database)
}
