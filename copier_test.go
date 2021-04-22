package copier_test

import (
	"database/sql"
	"errors"
	"fmt"
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
	Notes     []*string
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

	if len(employee.Notes) != len(user.Notes) {
		t.Fatalf("%v: Copy from slice doesn't work, employee notes len: %v, user: %v", testCase, len(employee.Notes), len(user.Notes))
	}

	for idx, note := range user.Notes {
		if note != *employee.Notes[idx] {
			t.Fatalf("%v: Copy from slice doesn't work, notes idx: %v employee: %v user: %v", testCase, idx, *employee.Notes[idx], note)
		}
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
	users := []User{{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, {Name: "Jinzhu2", Age: 22, Role: "Dev", Notes: []string{"hello world", "hello"}}}
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

func TestCopyFromSliceToSlice3(t *testing.T) {
	type CollectionAlias struct {
		CollectionName string `json:"collection_name"`
		Name           string `json:"name"`
	}

	expectedResult := []*CollectionAlias{
		{"collection", "collection_alias1"},
		{"collection", "collection_alias2"},
		{"collection", "collection_alias3"},
	}

	mockedResult := []*CollectionAlias{}
	copier.Copy(&mockedResult, &expectedResult)

	if len(mockedResult) != len(expectedResult) {
		t.Fatalf("failed to copy results")
	}

	for idx := range mockedResult {
		if mockedResult[idx].Name != mockedResult[idx].Name || mockedResult[idx].CollectionName != mockedResult[idx].CollectionName {
			t.Fatalf("failed to copy results")
		}
	}
}

func TestEmbeddedAndBase(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
		User       *User
	}

	type Embed struct {
		EmbedField1 int
		EmbedField2 int
		Base
	}

	base := Base{}
	embedded := Embed{}
	embedded.BaseField1 = 1
	embedded.BaseField2 = 2
	embedded.EmbedField1 = 3
	embedded.EmbedField2 = 4

	user := User{
		Name: "testName",
	}
	embedded.User = &user

	copier.Copy(&base, &embedded)

	if base.BaseField1 != 1 || base.User.Name != "testName" {
		t.Error("Embedded fields not copied")
	}

	base.BaseField1 = 11
	base.BaseField2 = 12
	user1 := User{
		Name: "testName1",
	}
	base.User = &user1

	copier.Copy(&embedded, &base)
	if embedded.BaseField1 != 11 || embedded.User.Name != "testName1" {
		t.Error("base fields not copied")
	}
}

func TestStructField(t *testing.T) {
	type Detail struct {
		Info1 string
		Info2 *string
	}

	type SimilarDetail struct {
		Info1 string
		Info2 *string
	}

	type UserWithDetailsPtr struct {
		Details []*Detail
		Detail  *Detail
	}
	type UserWithDetails struct {
		Details []Detail
		Detail  Detail
	}
	type UserWithSimilarDetailsPtr struct {
		Detail *SimilarDetail
	}
	type UserWithSimilarDetails struct {
		Detail SimilarDetail
	}
	type EmployeeWithDetails struct {
		Detail Detail
	}
	type EmployeeWithDetailsPtr struct {
		Detail *Detail
	}
	type EmployeeWithSimilarDetails struct {
		Detail SimilarDetail
	}
	type EmployeeWithSimilarDetailsPtr struct {
		Detail *SimilarDetail
	}

	optionsDeepCopy := copier.Option{
		DeepCopy: true,
	}

	checkDetail := func(t *testing.T, source Detail, target Detail) {
		if source.Info1 != target.Info1 {
			t.Errorf("info1 is diff: source: %v, target: %v", source.Info1, target.Info1)
		}

		if (source.Info2 != nil || target.Info2 != nil) && (*source.Info2 != *target.Info2) {
			t.Errorf("info2 is diff: source: %v, target: %v", *source.Info2, *target.Info2)
		}
	}

	t.Run("Should work without deepCopy", func(t *testing.T) {
		t.Run("Should work with same type and both ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{
				Detail:  &Detail{Info1: "hello", Info2: &info2},
				Details: []*Detail{{Info1: "hello", Info2: &info2}},
			}
			to := UserWithDetailsPtr{}
			copier.Copy(&to, from)

			checkDetail(t, *from.Detail, *to.Detail)

			*to.Detail.Info2 = "new value"
			if *from.Detail.Info2 != *to.Detail.Info2 {
				t.Fatalf("DeepCopy not enabled")
			}

			if len(to.Details) != len(to.Details) {
				t.Fatalf("slice should be copied")
			}

			for idx, detail := range from.Details {
				checkDetail(t, *detail, *to.Details[idx])
			}
		})

		t.Run("Should work with same type and both not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{
				Detail:  Detail{Info1: "hello", Info2: &info2},
				Details: []Detail{{Info1: "hello", Info2: &info2}},
			}
			to := UserWithDetails{}
			copier.Copy(&to, from)

			checkDetail(t, from.Detail, to.Detail)

			*to.Detail.Info2 = "new value"
			if *from.Detail.Info2 != *to.Detail.Info2 {
				t.Fatalf("DeepCopy not enabled")
			}

			if len(to.Details) != len(to.Details) {
				t.Fatalf("slice should be copied")
			}

			for idx, detail := range from.Details {
				checkDetail(t, detail, to.Details[idx])
			}
		})

		t.Run("Should work with different type and both ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{Detail: &Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetailsPtr{}
			copier.Copy(&to, from)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with different type and both not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{Detail: Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetails{}
			copier.Copy(&to, from)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with from ptr field and to not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{Detail: &Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetails{}
			copier.Copy(&to, from)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with from not ptr field and to ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{Detail: Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetailsPtr{}
			copier.Copy(&to, from)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})
	})

	t.Run("Should work with deepCopy", func(t *testing.T) {
		t.Run("Should work with same type and both ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{
				Detail:  &Detail{Info1: "hello", Info2: &info2},
				Details: []*Detail{{Info1: "hello", Info2: &info2}},
			}
			to := UserWithDetailsPtr{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			checkDetail(t, *from.Detail, *to.Detail)

			*to.Detail.Info2 = "new value"
			if *from.Detail.Info2 == *to.Detail.Info2 {
				t.Fatalf("DeepCopy enabled")
			}

			if len(to.Details) != len(to.Details) {
				t.Fatalf("slice should be copied")
			}

			for idx, detail := range from.Details {
				checkDetail(t, *detail, *to.Details[idx])
			}
		})
		t.Run("Should work with same type and both not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{
				Detail:  Detail{Info1: "hello", Info2: &info2},
				Details: []Detail{{Info1: "hello", Info2: &info2}},
			}
			to := UserWithDetails{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			checkDetail(t, from.Detail, to.Detail)

			*to.Detail.Info2 = "new value"
			if *from.Detail.Info2 == *to.Detail.Info2 {
				t.Fatalf("DeepCopy enabled")
			}

			if len(to.Details) != len(to.Details) {
				t.Fatalf("slice should be copied")
			}

			for idx, detail := range from.Details {
				checkDetail(t, detail, to.Details[idx])
			}
		})

		t.Run("Should work with different type and both ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{Detail: &Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetailsPtr{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with different type and both not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{Detail: Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetails{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with from ptr field and to not ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetailsPtr{Detail: &Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetails{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})

		t.Run("Should work with from not ptr field and to ptr field", func(t *testing.T) {
			info2 := "world"
			from := UserWithDetails{Detail: Detail{Info1: "hello", Info2: &info2}}
			to := EmployeeWithDetailsPtr{}
			copier.CopyWithOption(&to, from, optionsDeepCopy)

			newValue := "new value"
			to.Detail.Info2 = &newValue

			if to.Detail.Info1 == "" {
				t.Errorf("should not be empty")
			}
			if to.Detail.Info1 != from.Detail.Info1 {
				t.Errorf("should be the same")
			}
			if to.Detail.Info2 == from.Detail.Info2 {
				t.Errorf("should be different")
			}
		})
	})
}

func TestMapInterface(t *testing.T) {
	type Inner struct {
		IntPtr          *int
		unexportedField string
	}

	type Outer struct {
		Inner Inner
	}

	type DriverOptions struct {
		GenOptions map[string]interface{}
	}

	t.Run("Should work without deepCopy", func(t *testing.T) {
		intVal := 5
		outer := Outer{
			Inner: Inner{
				IntPtr:          &intVal,
				unexportedField: "hello",
			},
		}
		from := DriverOptions{
			GenOptions: map[string]interface{}{
				"key": outer,
			},
		}
		to := DriverOptions{}
		if err := copier.Copy(&to, &from); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		*to.GenOptions["key"].(Outer).Inner.IntPtr = 6

		if to.GenOptions["key"].(Outer).Inner.IntPtr != from.GenOptions["key"].(Outer).Inner.IntPtr {
			t.Errorf("should be the same")
		}
	})

	t.Run("Should work with deepCopy", func(t *testing.T) {
		intVal := 5
		outer := Outer{
			Inner: Inner{
				IntPtr:          &intVal,
				unexportedField: "Hello",
			},
		}
		from := DriverOptions{
			GenOptions: map[string]interface{}{
				"key": outer,
			},
		}
		to := DriverOptions{}
		if err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		}); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		*to.GenOptions["key"].(Outer).Inner.IntPtr = 6

		if to.GenOptions["key"].(Outer).Inner.IntPtr == from.GenOptions["key"].(Outer).Inner.IntPtr {
			t.Errorf("should be different")
		}
	})

	t.Run("Test copy map with nil interface", func(t *testing.T) {
		from := map[string]interface{}{"eventId": nil}
		to := map[string]interface{}{"eventId": nil}
		copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if v, ok := to["eventId"]; !ok || v != nil {
			t.Errorf("failed to deep copy map with nil, got %v", v)
		}

		from["eventId"] = 1
		if v, ok := to["eventId"]; !ok || v != nil {
			t.Errorf("failed to deep copy map with nil, got %v", v)
		}

		copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if v, ok := to["eventId"]; !ok || v != 1 {
			t.Errorf("failed to deep copy map with nil")
		}

		from["eventId"] = 2
		if v, ok := to["eventId"]; !ok || v != 1 {
			t.Errorf("failed to deep copy map with nil")
		}
	})

	t.Run("Test copy map with nested slice map", func(t *testing.T) {
		var out map[string]interface{}
		var value = map[string]interface{}{
			"list": []map[string]interface{}{
				{
					"shop_id": 123,
				},
			},
			"list2": []interface{}{
				map[string]interface{}{
					"shop_id": 123,
				},
			},
		}
		err := copier.CopyWithOption(&out, &value, copier.Option{IgnoreEmpty: false, DeepCopy: true})
		if err != nil {
			t.Fatalf("failed to deep copy nested map")
		}
		if fmt.Sprintf("%v", out) != fmt.Sprintf("%v", value) {
			t.Fatalf("failed to deep copy nested map")
		}
	})
}

func TestInterface(t *testing.T) {
	type Inner struct {
		IntPtr *int
	}

	type Outer struct {
		Inner Inner
	}

	type DriverOptions struct {
		GenOptions interface{}
	}

	t.Run("Should work without deepCopy", func(t *testing.T) {
		intVal := 5
		outer := Outer{
			Inner: Inner{
				IntPtr: &intVal,
			},
		}
		from := DriverOptions{
			GenOptions: outer,
		}
		to := DriverOptions{}
		if err := copier.Copy(&to, from); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		*to.GenOptions.(Outer).Inner.IntPtr = 6

		if to.GenOptions.(Outer).Inner.IntPtr != from.GenOptions.(Outer).Inner.IntPtr {
			t.Errorf("should be the same")
		}
	})

	t.Run("Should work with deepCopy", func(t *testing.T) {
		intVal := 5
		outer := Outer{
			Inner: Inner{
				IntPtr: &intVal,
			},
		}
		from := DriverOptions{
			GenOptions: outer,
		}
		to := DriverOptions{}
		if err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		}); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		*to.GenOptions.(Outer).Inner.IntPtr = 6

		if to.GenOptions.(Outer).Inner.IntPtr == from.GenOptions.(Outer).Inner.IntPtr {
			t.Errorf("should be different")
		}
	})
}

func TestSlice(t *testing.T) {
	type ElemOption struct {
		Value int
	}

	type A struct {
		X       []int
		Options []ElemOption
	}

	type B struct {
		X       []int
		Options []ElemOption
	}

	t.Run("Should work with simple slice", func(t *testing.T) {
		from := []int{1, 2}
		var to []int

		if err := copier.Copy(&to, from); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from[0] = 3
		from[1] = 4

		if to[0] == from[0] {
			t.Errorf("should be different")
		}

		if len(to) != len(from) {
			t.Errorf("should be the same length, got len(from): %v, len(to): %v", len(from), len(to))
		}
	})

	t.Run("Should work with empty slice", func(t *testing.T) {
		from := []int{}
		to := []int{}

		if err := copier.Copy(&to, from); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if to == nil {
			t.Errorf("should be not nil")
		}
	})

	t.Run("Should work without deepCopy", func(t *testing.T) {
		x := []int{1, 2}
		options := []ElemOption{
			{Value: 10},
			{Value: 20},
		}
		from := A{X: x, Options: options}
		to := B{}

		if err := copier.Copy(&to, from); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.X[0] = 3
		from.X[1] = 4
		from.Options[0].Value = 30
		from.Options[1].Value = 40

		if to.X[0] != from.X[0] {
			t.Errorf("should be the same")
		}

		if len(to.X) != len(from.X) {
			t.Errorf("should be the same length, got len(from.X): %v, len(to.X): %v", len(from.X), len(to.X))
		}

		if to.Options[0].Value != from.Options[0].Value {
			t.Errorf("should be the same")
		}

		if to.Options[0].Value != from.Options[0].Value {
			t.Errorf("should be the same")
		}

		if len(to.Options) != len(from.Options) {
			t.Errorf("should be the same")
		}
	})

	t.Run("Should work with deepCopy", func(t *testing.T) {
		x := []int{1, 2}
		options := []ElemOption{
			{Value: 10},
			{Value: 20},
		}
		from := A{X: x, Options: options}
		to := B{}

		if err := copier.CopyWithOption(&to, from, copier.Option{
			DeepCopy: true,
		}); nil != err {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.X[0] = 3
		from.X[1] = 4
		from.Options[0].Value = 30
		from.Options[1].Value = 40

		if to.X[0] == from.X[0] {
			t.Errorf("should be different")
		}

		if len(to.X) != len(from.X) {
			t.Errorf("should be the same length, got len(from.X): %v, len(to.X): %v", len(from.X), len(to.X))
		}

		if to.Options[0].Value == from.Options[0].Value {
			t.Errorf("should be different")
		}

		if len(to.Options) != len(from.Options) {
			t.Errorf("should be the same")
		}
	})
}

func TestAnonymousFields(t *testing.T) {
	t.Run("Should work with unexported ptr fields", func(t *testing.T) {
		type nested struct {
			A string
		}
		type parentA struct {
			*nested
		}
		type parentB struct {
			*nested
		}

		from := parentA{nested: &nested{A: "a"}}
		to := parentB{}

		err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.nested.A = "b"

		if to.nested != nil {
			t.Errorf("should be nil")
		}
	})
	t.Run("Should work with unexported fields", func(t *testing.T) {
		type nested struct {
			A string
		}
		type parentA struct {
			nested
		}
		type parentB struct {
			nested
		}

		from := parentA{nested: nested{A: "a"}}
		to := parentB{}

		err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.nested.A = "b"

		if to.nested.A == from.nested.A {
			t.Errorf("should be different")
		}
	})

	t.Run("Should work with exported ptr fields", func(t *testing.T) {
		type Nested struct {
			A string
		}
		type parentA struct {
			*Nested
		}
		type parentB struct {
			*Nested
		}

		fieldValue := "a"
		from := parentA{Nested: &Nested{A: fieldValue}}
		to := parentB{}

		err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.Nested.A = "b"

		if to.Nested.A != fieldValue {
			t.Errorf("should not change")
		}
	})

	t.Run("Should work with exported fields", func(t *testing.T) {
		type Nested struct {
			A string
		}
		type parentA struct {
			Nested
		}
		type parentB struct {
			Nested
		}

		fieldValue := "a"
		from := parentA{Nested: Nested{A: fieldValue}}
		to := parentB{}

		err := copier.CopyWithOption(&to, &from, copier.Option{
			DeepCopy: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		from.Nested.A = "b"

		if to.Nested.A != fieldValue {
			t.Errorf("should not change")
		}
	})
}

type someStruct struct {
	IntField  int
	UIntField uint64
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
	D string
	E *someStruct
}

type structSameName2 struct {
	A string
	B time.Time
	C int64
	D string
	E *someStruct
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

type Foo1 struct {
	Name string
	Age  int32
}

type Foo2 struct {
	Name string
}

type StructWithMap1 struct {
	Map map[int]Foo1
}

type StructWithMap2 struct {
	Map map[int32]Foo2
}

func TestCopyMapOfStruct(t *testing.T) {
	obj1 := StructWithMap1{Map: map[int]Foo1{2: {Name: "A pure foo"}}}
	obj2 := &StructWithMap2{}
	err := copier.Copy(obj2, obj1)
	if err != nil {
		t.Error("Should not raise error")
	}
	for k, v1 := range obj1.Map {
		v2, ok := obj2.Map[int32(k)]
		if !ok || v1.Name != v2.Name {
			t.Errorf("Map should be copied")
		}
	}
}

func TestCopyMapOfInt(t *testing.T) {
	map1 := map[int]int{3: 6, 4: 8}
	map2 := map[int32]int8{}
	err := copier.Copy(&map2, map1)
	if err != nil {
		t.Error("Should not raise error")
	}

	for k, v1 := range map1 {
		v2, ok := map2[int32(k)]
		if !ok || v1 != int(v2) {
			t.Errorf("Map should be copied")
		}
	}
}

func TestCopyWithOption(t *testing.T) {
	from := structSameName2{D: "456", E: &someStruct{IntField: 100, UIntField: 1000}}
	to := &structSameName1{A: "123", B: 2, C: time.Now(), D: "123", E: &someStruct{UIntField: 5000}}
	if err := copier.CopyWithOption(to, &from, copier.Option{IgnoreEmpty: true}); err != nil {
		t.Error("Should not raise error")
	}

	if to.A == from.A {
		t.Errorf("Field A should not be copied")
	} else if to.D != from.D {
		t.Errorf("Field D should be copied")
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

func TestScanFromPtrToSqlNullable(t *testing.T) {

	var (
		from struct {
			S    string
			Sptr *string
			T1   sql.NullTime
			T2   sql.NullTime
			T3   *time.Time
		}

		to struct {
			S    sql.NullString
			Sptr sql.NullString
			T1   time.Time
			T2   *time.Time
			T3   sql.NullTime
		}

		s string

		err error
	)

	s = "test"
	from.S = s
	from.Sptr = &s

	if from.T1.Valid || from.T2.Valid {
		t.Errorf("Must be not valid")
	}

	err = copier.Copy(&to, from)
	if err != nil {
		t.Error("Should not raise error")
	}

	if !to.T1.IsZero() {
		t.Errorf("to.T1 should be Zero but %v", to.T1)
	}

	if to.T2 != nil && !to.T2.IsZero() {
		t.Errorf("to.T2 should be Zero but %v", to.T2)
	}

	now := time.Now()

	from.T1.Scan(now)
	from.T2.Scan(now)

	err = copier.Copy(&to, from)
	if err != nil {
		t.Error("Should not raise error")
	}

	if to.S.String != from.S {
		t.Errorf("Field S should be copied")
	}

	if to.Sptr.String != *from.Sptr {
		t.Errorf("Field Sptr should be copied")
	}

	if from.T1.Time != to.T1 {
		t.Errorf("Fields T1 fields should be equal")
	}

	if from.T2.Time != *to.T2 {
		t.Errorf("Fields T2 fields should be equal")
	}
}

func TestDeepCopyInterface(t *testing.T) {
	var m = make(map[string]string)
	m["a"] = "ccc"

	from := []interface{}{[]int{7, 8, 9}, 2, 3, m, errors.New("aaaa")}
	var to []interface{}

	copier.CopyWithOption(&to, &from, copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    true,
	})

	from[0].([]int)[0] = 10
	from[1] = "3"
	from[3].(map[string]string)["a"] = "bbb"

	if fmt.Sprint(to[0]) != fmt.Sprint([]int{7, 8, 9}) {
		t.Errorf("to value failed to be deep copied")
	}

	if fmt.Sprint(to[1]) != "2" {
		t.Errorf("to value failed to be deep copied")
	}

	if to[3].(map[string]string)["a"] != "ccc" {
		t.Errorf("to value failed to be deep copied")
	}
}
