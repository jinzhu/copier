package copier

import (
	"encoding/json"

	"reflect"

	"testing"
)

type User struct {
	Name  string
	Role  string
	Age   int32
	Notes []string
}

func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Age       int32
	EmployeId int64
	DoubleAge int32
	SuperRule string
	Notes     []string
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func TestCopyStruct(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	employee := Employee{}

	Copy(&employee, &user)

	if employee.Name != "Jinzhu" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employee.Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employee.DoubleAge != 36 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employee.SuperRule != "Super Admin" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Copy a map")
	}

	user.Notes = append(user.Notes, "welcome")
	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("User's Note should be changed")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Employee's Note should not be changed")
	}

	employee.Notes = append(employee.Notes, "golang")
	if !reflect.DeepEqual(employee.Notes, []string{"hello world", "golang"}) {
		t.Errorf("Employee's Note should be changed")
	}

	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("Employee's Note should not be changed")
	}
}

func TestCopySlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	users := []User{{Name: "jinzhu 2", Age: 30, Role: "Dev"}}
	employees := []Employee{}

	Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	Copy(&employees, &users)
	if len(employees) != 2 {
		t.Errorf("Should have two elems when copy additional slice to slice")
	}

	if employees[0].Name != "Jinzhu" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[0].Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[0].DoubleAge != 36 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[0].SuperRule != "Super Admin" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	if employees[1].Name != "jinzhu 2" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[1].Age != 30 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[1].DoubleAge != 60 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[1].SuperRule != "Super Dev" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	employee := employees[0]
	user.Notes = append(user.Notes, "welcome")
	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("User's Note should be changed")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Employee's Note should not be changed")
	}

	employee.Notes = append(employee.Notes, "golang")
	if !reflect.DeepEqual(employee.Notes, []string{"hello world", "golang"}) {
		t.Errorf("Employee's Note should be changed")
	}

	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("Employee's Note should not be changed")
	}
}

func BenchmarkCopyStruct(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		Copy(&Employee{}, &user)
	}
}

func BenchmarkNamaCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			Name:      user.Name,
			Age:       user.Age,
			DoubleAge: user.DoubleAge(),
			Notes:     user.Notes,
		}
		employee.Role(user.Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)
		employee.DoubleAge = user.DoubleAge()
		employee.Role(user.Role)
	}
}
