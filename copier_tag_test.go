package copier_test

import (
	"copier"
	"errors"
	"reflect"
	"testing"
	"time"


)

type UserTag struct {
	Name     string     `mson:"Name_string"`
	Birthday *time.Time `mson:"Birthday_time_Time"`
	Nickname string     `mson:"Nickname_string"`
	Role     string     `mson:"Role_string"`
	Age      int32      `mson:"Age_int32"`
	FakeAge  *int32     `mson:"FakeAge_int32"`
	Notes    []string   `mson:"Notes_string_slice"`
	flags    []byte     `mson:"flags_byte_slice"`
}

func (user UserTag) DoubleAge_int32() int32 {
	return 2 * user.Age
}

type EmployeeTag struct {
	Name      string     `mson:"Name_string"`
	Birthday  *time.Time `mson:"Birthday_time_Time"`
	Nickname  *string    `mson:"Nickname_string"`
	Age       int64      `mson:"Age_int64"`
	FakeAge   int        `mson:"FakeAge_int"`
	EmployeeID int64     `mson:"EmployeeId_int64"`
	DoubleAge int32      `mson:"DoubleAge_int32"`
	SuperRule string     `mson:"SuperRule_string"`
	Notes     []string   `mson:"Notes_string_slice"`
	flags     []byte     `mson:"flags_byte_slice"`
}

func (employee *EmployeeTag) FakeAge_int32(FakeAge *int32) {
	if FakeAge == nil{
		return
	}
	employee.FakeAge = int(*FakeAge)
}

func (employee *EmployeeTag) Age_int32(age int32) {
	employee.Age = int64(age)
}

func (employee *EmployeeTag) Role_string(role string) {
	employee.SuperRule = "Super " + role
}

func checkEmployeeTag(employee EmployeeTag, user UserTag, t *testing.T, testCase string) {
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
	if employee.DoubleAge != user.DoubleAge_int32() {
		t.Errorf("%v: Copy from method doesn't work", testCase)
	}
	if employee.SuperRule != "Super "+user.Role {
		t.Errorf("%v: Copy to method doesn't work", testCase)
	}
	if !reflect.DeepEqual(employee.Notes, user.Notes) {
		t.Errorf("%v: Copy from slice doen't work", testCase)
	}
}

func TestCopySameStructWithPointerFieldTag(t *testing.T) {
	var fakeAge int32 = 12
	var currentTime time.Time = time.Now()
	user := &UserTag{Birthday: &currentTime, Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge,
		Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	newUser := &UserTag{}
	copier.CopyByTag(newUser, user,"mson")
	if user.Birthday == newUser.Birthday {
		t.Errorf("TestCopySameStructWithPointerField: copy Birthday failed since they need to have different address")
	}

	if user.FakeAge == newUser.FakeAge {
		t.Errorf("TestCopySameStructWithPointerField: copy FakeAge failed since they need to have different address")
	}
}

func checkEmployeeTag2(employee EmployeeTag, user *UserTag, t *testing.T, testCase string) {
	if user == nil {
		if employee.Name != "" || employee.Nickname != nil || employee.Birthday != nil || employee.Age != 0 ||
			employee.DoubleAge != 0 || employee.FakeAge != 0 || employee.SuperRule != "" || employee.Notes != nil {
			t.Errorf("%v : employee should be empty", testCase)
		}
		return
	}

	checkEmployeeTag(employee, *user, t, testCase)
}

func TestCopyStructTag(t *testing.T) {
	var fakeAge int32 = 12
	user := UserTag{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin",
		Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	employee := EmployeeTag{}

	if err := copier.CopyByTag(employee, &user,"mson"); err == nil {
		t.Errorf("Copy to unaddressable value should get error")
	}

	copier.CopyByTag(&employee, &user,"mson")
	checkEmployeeTag(employee, user, t, "Copy From Ptr To Ptr")

	employee2 := EmployeeTag{}
	copier.CopyByTag(&employee2, user,"mson")
	checkEmployeeTag(employee2, user, t, "Copy From Struct To Ptr")

	employee3 := EmployeeTag{}
	ptrToUser := &user
	copier.CopyByTag(&employee3, &ptrToUser,"mson")
	checkEmployeeTag(employee3, user, t, "Copy From Double Ptr To Ptr")

	employee4 := &EmployeeTag{}
	copier.CopyByTag(&employee4, user,"mson")
	checkEmployeeTag(*employee4, user, t, "Copy From Ptr To Double Ptr")
}

func TestCopyFromStructToSliceTag(t *testing.T) {
	user := UserTag{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	employees := []EmployeeTag{}

	if err := copier.CopyByTag(employees, &user,"mson"); err != nil && len(employees) != 0 {
		t.Errorf("Copy to unaddressable value should get error")
	}

	if copier.CopyByTag(&employees, &user,"mson"); len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployeeTag(employees[0], user, t, "Copy From Struct To Slice Ptr")
	}

	employees2 := &[]EmployeeTag{}
	if copier.CopyByTag(&employees2, user,"mson"); len(*employees2) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployeeTag((*employees2)[0], user, t, "Copy From Struct To Double Slice Ptr")
	}

	employees3 := []*EmployeeTag{}
	if copier.CopyByTag(&employees3, user,"mson"); len(employees3) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployeeTag(*(employees3[0]), user, t, "Copy From Struct To Ptr Slice Ptr")
	}

	employees4 := &[]*EmployeeTag{}
	if copier.CopyByTag(&employees4, user,"mson"); len(*employees4) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployeeTag(*((*employees4)[0]), user, t, "Copy From Struct To Double Ptr Slice Ptr")
	}
}

func TestCopyFromSliceToSliceTag(t *testing.T) {
	users := []UserTag{UserTag{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}},
		UserTag{Name: "Jinzhu2", Age: 22, Role: "Dev", Notes: []string{"hello world", "hello"}}}
	employees := []EmployeeTag{}

	if copier.CopyByTag(&employees, users,"mson"); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployeeTag(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]EmployeeTag{}
	if copier.CopyByTag(&employees2, &users,"mson"); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployeeTag((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*EmployeeTag{}
	if copier.CopyByTag(&employees3, users,"mson"); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployeeTag(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*EmployeeTag{}
	if copier.CopyByTag(&employees4, users,"mson"); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployeeTag(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestCopyFromSliceToSliceTag2(t *testing.T) {
	users := []*UserTag{{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, nil}
	employees := []EmployeeTag{}

	if copier.CopyByTag(&employees, users,"mson"); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag2(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployeeTag2(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]EmployeeTag{}
	if copier.CopyByTag(&employees2, &users,"mson"); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag2((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployeeTag2((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*EmployeeTag{}
	if copier.CopyByTag(&employees3, users,"mson"); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag2(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployeeTag2(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*EmployeeTag{}
	if copier.CopyByTag(&employees4, users,"mson"); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployeeTag2(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployeeTag2(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestEmbeddedAndBaseTag(t *testing.T) {
	type Base struct {
		BaseField1 int `mson:"BaseField1_int"`
		BaseField2 int `mson:"baseField2_int"`
		User *UserTag     `mson:"User_User"`
	}

	type Embed struct {
		EmbedField1 int `mson:"EmbedField1_int"`
		EmbedField2 int `mson:"EmbedField2_int"`
		Base
	}

	base := Base{}
	embeded := Embed{}
	embeded.BaseField1 = 1
	embeded.BaseField2 = 2
	embeded.EmbedField1 = 3
	embeded.EmbedField2 = 4

	user:=UserTag{
		Name:"testName",
	}
	embeded.User=&user

	copier.CopyByTag(&base, &embeded,"mson")

	if base.BaseField1 != 1 || base.User.Name!="testName"{
		t.Error("Embedded fields not copied")
	}

	base.BaseField1=11
	base.BaseField2=12
	user1:=UserTag{
		Name:"testName1",
	}
	base.User=&user1

	copier.CopyByTag(&embeded,&base,"mson")

	if embeded.BaseField1 != 11 || embeded.User.Name!="testName1" {
		t.Error("base fields not copied")
	}
}

type structSameName1Tag struct {
	A string    `mson:"A_string"`
	B int64     `mson:"B_int64"`
	C time.Time `mson:"C_time_Time"`
}

type structSameName2Tag struct {
	A string    `mson:"A_string"`
	B time.Time `mson:"B_time_Time"`
	C int64     `mson:"C_int64"`
}

func (s *structSameName2Tag) C_time_Time(t time.Time){
	s.C = t.Unix()
}
func (s *structSameName2Tag) B_int64(B int64){
	s.B = time.Unix(B,0)
}
func TestCopyFieldsWithSameNameButDifferentTypesTag(t *testing.T) {
	obj1 := structSameName1Tag{A: "123", B: 2, C: time.Now()}
	obj2 := &structSameName2Tag{}
	err := copier.CopyByTag(obj2, &obj1,"mson")
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}
	if obj2.B.Unix() != obj1.B{
		t.Errorf("Field B should be copied")
	}
	if obj2.C != obj1.C.Unix(){
		t.Errorf("Field C should be copied")
	}
}

type ScannerValueTag struct {
	V int  `json:"mson:"V_int"`
}

func (s *ScannerValue) ScanTag(src interface{}) error {
	return errors.New("I failed")
}

type ScannerStructTag struct {
	V *ScannerValueTag   `mson:"V_ScannerValueTag"`
}

type ScannerStructToTag struct {
	V *ScannerValueTag `mson:"V_ScannerValueTag"`
}

func TestScannerTag(t *testing.T) {
	s := &ScannerStructTag{
		V: &ScannerValueTag{
			V: 12,
		},
	}

	s2 := &ScannerStructToTag{}

	err := copier.CopyByTag(s2, s,"mson")
	if err != nil {
		t.Error("Should not raise error")
	}

	if s.V.V != s2.V.V {
		t.Errorf("Field V should be copied")
	}
}
