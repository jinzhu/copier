package copier_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

type User struct {
	Name     string
	Birthday *time.Time
	Nickname string
	Role     string
	Age      int32
	FakeAge  *int32
	Notes    []string
	flags    []byte
}

func (user User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Birthday  *time.Time
	Nickname  *string
	Age       int64
	FakeAge   int
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
	if employee.Nickname == nil || *employee.Nickname != user.Nickname {
		t.Errorf("%v: NickName haven't been copied correctly.", testCase)
	}
	if employee.Birthday == nil && user.Birthday != nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Birthday != nil && user.Birthday == nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Birthday != nil && user.Birthday != nil &&
		!employee.Birthday.Equal(*(user.Birthday)) {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Age != int64(user.Age) {
		t.Errorf("%v: Age haven't been copied correctly.", testCase)
	}
	if user.FakeAge != nil && employee.FakeAge != int(*user.FakeAge) {
		t.Errorf("%v: FakeAge haven't been copied correctly.", testCase)
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

func TestCopySameStructWithPointerField(t *testing.T) {
	var fakeAge int32 = 12
	var currentTime time.Time = time.Now()
	user := &User{Birthday: &currentTime, Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	newUser := &User{}
	copier.Copy(newUser, user)
	if user.Birthday == newUser.Birthday {
		t.Errorf("TestCopySameStructWithPointerField: copy Birthday failed since they need to have different address")
	}

	if user.FakeAge == newUser.FakeAge {
		t.Errorf("TestCopySameStructWithPointerField: copy FakeAge failed since they need to have different address")
	}
}

func checkEmployee2(employee Employee, user *User, t *testing.T, testCase string) {
	if user == nil {
		if employee.Name != "" || employee.Nickname != nil || employee.Birthday != nil || employee.Age != 0 ||
			employee.DoubleAge != 0 || employee.FakeAge != 0 || employee.SuperRule != "" || employee.Notes != nil {
			t.Errorf("%v : employee should be empty", testCase)
		}
		return
	}

	checkEmployee(employee, *user, t, testCase)
}

func TestCopyStruct(t *testing.T) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	employee := Employee{}

	if err := copier.Copy(employee, &user); err == nil {
		t.Errorf("Copy to unaddressable value should get error")
	}

	copier.Copy(&employee, &user)
	checkEmployee(employee, user, t, "Copy From Ptr To Ptr")

	employee2 := Employee{}
	copier.Copy(&employee2, user)
	checkEmployee(employee2, user, t, "Copy From Struct To Ptr")

	employee3 := Employee{}
	ptrToUser := &user
	copier.Copy(&employee3, &ptrToUser)
	checkEmployee(employee3, user, t, "Copy From Double Ptr To Ptr")

	employee4 := &Employee{}
	copier.Copy(&employee4, user)
	checkEmployee(*employee4, user, t, "Copy From Ptr To Double Ptr")
}

func TestCopyFromStructToSlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	employees := []Employee{}

	if err := copier.Copy(employees, &user); err != nil && len(employees) != 0 {
		t.Errorf("Copy to unaddressable value should get error")
	}

	if copier.Copy(&employees, &user); len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(employees[0], user, t, "Copy From Struct To Slice Ptr")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, user); len(*employees2) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee((*employees2)[0], user, t, "Copy From Struct To Double Slice Ptr")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, user); len(employees3) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*(employees3[0]), user, t, "Copy From Struct To Ptr Slice Ptr")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, user); len(*employees4) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*((*employees4)[0]), user, t, "Copy From Struct To Double Ptr Slice Ptr")
	}
}

func TestCopyFromSliceToSlice(t *testing.T) {
	users := []User{User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, User{Name: "Jinzhu2", Age: 22, Role: "Dev", Notes: []string{"hello world", "hello"}}}
	employees := []Employee{}

	if copier.Copy(&employees, users); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployee(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, &users); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployee((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, users); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployee(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, users); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployee(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestCopyFromSliceToSlice2(t *testing.T) {
	users := []*User{{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, nil}
	employees := []Employee{}

	if copier.Copy(&employees, users); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee2(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployee2(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, &users); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee2((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployee2((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, users); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee2(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployee2(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, users); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee2(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployee2(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestEmbeddedAndBase(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
		User *User
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

	user:=User{
		Name:"testName",
	}
	embeded.User=&user

	copier.Copy(&base, &embeded)

	if base.BaseField1 != 1 || base.User.Name!="testName"{
		t.Error("Embedded fields not copied")
	}

	base.BaseField1=11
	base.BaseField2=12
	user1:=User{
		Name:"testName1",
	}
	base.User=&user1

	copier.Copy(&embeded,&base)

	if embeded.BaseField1 != 11 || embeded.User.Name!="testName1" {
		t.Error("base fields not copied")
	}
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
}

type structSameName2 struct {
	A string
	B time.Time
	C int64
}

func TestCopyFieldsWithSameNameButDifferentTypes(t *testing.T) {
	obj1 := structSameName1{A: "123", B: 2, C: time.Now()}
	obj2 := &structSameName2{}
	err := copier.Copy(obj2, &obj1)
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}
}

type ScannerValue struct {
	V int
}

func (s *ScannerValue) Scan(src interface{}) error {
	return errors.New("I failed")
}

type ScannerStruct struct {
	V *ScannerValue
}

type ScannerStructTo struct {
	V *ScannerValue
}

func TestScanner(t *testing.T) {
	s := &ScannerStruct{
		V: &ScannerValue{
			V: 12,
		},
	}

	s2 := &ScannerStructTo{}

	err := copier.Copy(s2, s)
	if err != nil {
		t.Error("Should not raise error")
	}

	if s.V.V != s2.V.V {
		t.Errorf("Field V should be copied")
	}
}
