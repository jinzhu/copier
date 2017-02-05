package copier_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
)

type User struct {
	Name  string
	Role  string
	Age   int32
	Notes []string
	flags []byte
}

func (user User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Age       int32
	EmployeID int64
	DoubleAge int32
	SuperRule string
	Notes     []string
	flags     []byte
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func checkEmployee(employee Employee, user User, t *testing.T, testCase string) {
	if employee.Name != user.Name {
		t.Errorf("%v: Name haven't been copied correctly.", testCase)
	}
	if employee.Age != user.Age {
		t.Errorf("%v: Age haven't been copied correctly.", testCase)
	}
	if employee.DoubleAge != user.DoubleAge() {
		t.Errorf("%v: Copy from method doesn't work", testCase)
	}
	if employee.SuperRule != "Super "+user.Role {
		t.Errorf("%v: Copy to method doesn't work", testCase)
	}
	if !reflect.DeepEqual(employee.Notes, user.Notes) {
		t.Errorf("%v: Copy from slice doen't work", testCase)
	}
}

func TestCopyStruct(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	employee := Employee{}

	copier.Copy(&employee, &user)
	checkEmployee(employee, user, t, "Copy From Ptr To Ptr")

	employee2 := Employee{}
	copier.Copy(&employee2, user)
	checkEmployee(employee2, user, t, "Copy From Struct To Ptr")
}

func TestCopySlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	users := []User{{Name: "jinzhu 2", Age: 30, Role: "Dev"}}
	employees := []Employee{}

	copier.Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	copier.Copy(&employees, &users)
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

func TestCopySliceWithPtr(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	user2 := &User{Name: "jinzhu 2", Age: 30, Role: "Dev"}
	users := []*User{user2}
	employees := []*Employee{}

	copier.Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	copier.Copy(&employees, &users)
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

func TestEmbedded(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
	}

	type Embed struct {
		EmbedField1 int
		EmbedField2 int
		Base
	}

	base := Base{}
	embeded := Embed{}
	embeded.BaseField1 = 1
	embeded.BaseField2 = 2
	embeded.EmbedField1 = 3
	embeded.EmbedField2 = 4

	copier.Copy(&base, &embeded)

	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
	}
}

type TypeStruct1 struct {
	Field1 string
	Field2 string
	Field3 TypeStruct2
	Field4 *TypeStruct2
	Field5 []*TypeStruct2
}

type TypeStruct2 struct {
	Field1 int
	Field2 string
}

type TypeStruct3 struct {
	Field1 interface{}
	Field2 string
	Field3 TypeStruct4
	Field4 *TypeStruct4
	Field5 []*TypeStruct4
}

type TypeStruct4 struct {
	field1 int
	Field2 string
}

type TypeStruct5 struct {
	field1 string
	Field2 string
}

func (t *TypeStruct4) Field1(i int) {
	t.field1 = i
}

func (t *TypeStruct5) Field1(i interface{}) {
	if v, ok := i.(string); ok {
		t.field1 = v
	}
}

func TestDifferentType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The copy did panic")
		}
	}()

	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}

	ts2 := &TypeStruct2{}

	copier.Copy(ts2, ts)
}

func TestDifferentTypeMethod(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The copy did panic")
		}
	}()

	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}

	ts4 := &TypeStruct4{}

	copier.Copy(ts4, ts)
}

func TestAssignableType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The copy did panic: %v", r)
		}
	}()

	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
		Field3: TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field4: &TypeStruct2{
			Field1: 666,
			Field2: "str2",
		},
		Field5: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
	}

	ts3 := &TypeStruct3{}

	copier.Copy(&ts3, &ts)

	if v, ok := ts3.Field1.(string); !ok {
		t.Error("Assign to interface{} type did not succeed")
	} else if v != "str1" {
		t.Error("String haven't been copied correctly")
	}

	if ts3.Field4 == nil {
		t.Error("nil Field4")
	} else if ts3.Field4.Field2 != ts.Field4.Field2 {
		t.Errorf("Field4 differs %v", ts3.Field4)
	}
}

func TestPointerArray(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The copy did panic: %v", r)
		}
	}()

	ts := []*TypeStruct1{
		{
			Field1: "str1",
			Field2: "str2",
			Field3: TypeStruct2{

				Field1: 666,
				Field2: "str2",
			},
			Field4: &TypeStruct2{

				Field1: 666,
				Field2: "str2",
			},
			Field5: []*TypeStruct2{
				{
					Field1: 666,
					Field2: "str2",
				},
			},
		},
	}

	ts3 := []*TypeStruct3{}

	copier.Copy(&ts3, &ts)

	for i := range ts {
		if v, ok := ts3[i].Field1.(string); !ok {
			t.Error("Assign to interface{} type did not succeed")
		} else if v != "str1" {
			t.Error("String haven't been copied correctly")
		}

		if ts3[i].Field2 != ts[i].Field2 {
			t.Error("String haven't been copied correctly")
		}

		if ts3[i].Field3.Field2 != ts[i].Field3.Field2 {
			t.Errorf("String haven't been copied correctly %+v vs %+v", ts3[i].Field3, ts[i].Field3)
		}

		if ts3[i].Field4 == nil {
			t.Error("nil Field4")
		} else if ts3[i].Field4.Field2 != ts[i].Field4.Field2 {
			t.Errorf("Field4 differs %v", ts3[i].Field4)
		}

		if len(ts3[i].Field5) != len(ts[i].Field5) {
			t.Errorf("Field5 size differs %v and %v", len(ts3[i].Field5), len(ts[i].Field5))
		}
	}
}

func TestAssignableTypeMethod(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The copy did panic")
		}
	}()

	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}

	ts5 := &TypeStruct5{}

	copier.Copy(ts5, ts)

	if ts5.field1 != "str1" {
		t.Error("String haven't been copied correctly through method")
	}
}
func BenchmarkCopyStruct(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		copier.Copy(&Employee{}, &user)
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
