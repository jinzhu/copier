package copier

import (
	"encoding/json"

	"reflect"
	"unsafe"

	"testing"
)

type User struct {
	StringValue4 string
	Role         string
	Age          int32
	Notes        []string
	flags        []byte
}

func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	StringValue4 string
	Age          int32
	EmployeID    int64
	DoubleAge    int32
	SuperRule    string
	Notes        []string
	flags        []byte
}

type Base struct {
	BaseField1 int
	BaseField2 int
}

type HaveEmbed struct {
	EmbedField1 int
	EmbedField2 int
	Base
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

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func TestEmbedded(t *testing.T) {
	defer panicHandler(t)
	base := Base{}
	embeded := HaveEmbed{}
	embeded.BaseField1 = 1
	embeded.BaseField2 = 2
	embeded.EmbedField1 = 3
	embeded.EmbedField2 = 4

	Copy(&base, &embeded)

	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
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

	Copy(ts2, ts)
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

	Copy(ts4, ts)
}

func TestAssignableType(t *testing.T) {
	defer panicHandler(t)

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

	err := Copy(&ts3, &ts)

	if err != nil {
		t.Errorf("copying error: %+v", err)
	}

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
	defer panicHandler(t)

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

	err := Copy(&ts3, &ts)

	if nil != err {
		t.Errorf("copy error %v", err)
	}

	if len(ts3) != len(ts) {
		t.Errorf("Arrays of different length: original %d and destination %d", len(ts), len(ts3))
	}

	for i := range ts3 {
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

func TestArray(t *testing.T) {
	defer panicHandler(t)

	ts := []TypeStruct1{
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

	ts3 := []TypeStruct3{}

	Copy(&ts3, &ts)

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
	defer panicHandler(t)

	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}

	ts5 := &TypeStruct5{}

	Copy(ts5, ts)

	if ts5.field1 != "str1" {
		t.Error("String haven't been copied correctly through method")
	}

	if ts5.Field2 != "str2" {
		t.Error("String haven't been copied correctly through method")
	}
}

func TestCopyStruct(t *testing.T) {
	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}, flags: []byte{'x'}}
	employee := Employee{}

	Copy(&employee, &user)

	if employee.StringValue4 != "Jinzhu" {
		t.Errorf("StringValue4 haven't been copied correctly.")
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
	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	users := []User{{StringValue4: "jinzhu 2", Age: 30, Role: "Dev"}}
	employees := []Employee{}

	Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice, got %d", len(employees))
	}

	Copy(&employees, &users)
	if len(employees) != 2 {
		t.Errorf("Should have two elems when copy additional slice to slice")
	}

	if employees[0].StringValue4 != "Jinzhu" {
		t.Errorf("StringValue4 haven't been copied correctly.")
	}
	if employees[0].Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[0].DoubleAge != 36 {
		t.Errorf("Copy from method doesn't work %v", employees[0].DoubleAge)
	}
	if employees[0].SuperRule != "Super Admin" {
		t.Errorf("Copy should support copy to method")
	}

	if employees[1].StringValue4 != "jinzhu 2" {
		t.Errorf("StringValue4 haven't been copied correctly.")
	}
	if employees[1].Age != 30 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[1].DoubleAge != 60 {
		t.Errorf("Copy from method doesn't work")
	}
	if employees[1].SuperRule != "Super Dev" {
		t.Errorf("Copy should support copy to method")
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
	defer panicHandler(t)

	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	user2 := &User{StringValue4: "jinzhu 2", Age: 30, Role: "Dev"}
	users := []*User{user2}
	employees := []*Employee{}

	Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	Copy(&employees, &users)
	if len(employees) != 2 {
		t.Errorf("Should have two elems when copy additional slice to slice")
	}

	if employees[0].StringValue4 != "Jinzhu" {
		t.Errorf("StringValue4 haven't been copied correctly.")
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

	if employees[1].StringValue4 != "jinzhu 2" {
		t.Errorf("StringValue4 haven't been copied correctly.")
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
	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		Copy(&Employee{}, &user)
	}
}

func BenchmarkNamaCopy(b *testing.B) {
	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			StringValue4: user.StringValue4,
			Age:          user.Age,
			DoubleAge:    user.DoubleAge(),
			Notes:        user.Notes,
		}
		employee.Role(user.Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	user := User{StringValue4: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)
		employee.DoubleAge = user.DoubleAge()
		employee.Role(user.Role)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ComplicatedStruct1 struct {
	FloatValue1           *float64
	FloatValue2           *float64
	ComplicatedStruct2Ptr *ComplicatedStruct2
	ComplicatedArray1Val  ComplicatedArray1
}

type ComplicatedStruct2 struct {
	FloatArrayVal2 [][]float64
	StringValue1   string
}

type ComplicatedArray1 []*ComplicatedStruct3

type ComplicatedStruct3 struct {
	FloatValue1          float64
	FloatValue2          float64
	ComplicatedArray2Val ComplicatedArray2
	StringValue2         string
}

type ComplicatedArray2 []*ComplicatedStruct4

type ComplicatedStruct4 struct {
	FloatValue1           float64
	FloatValue2           float64
	ComplicatedStruct2Ptr *ComplicatedStruct2
	ComplicatedStruct5Ptr *ComplicatedStruct5
	StringValue3          string
	StringValue4          string
}

type ComplicatedStruct5 struct {
	StringValue5   float64
	StringValue6   float64
	FloatArrayVal1 []float64
	StringValue7   string
	StringValue1   string
}

func TestPointerArrayComplex(t *testing.T) {
	defer panicHandler(t)

	FloatValue1 := 13.666
	FloatValue2 := 666.13

	ComplicatedStruct2Ptr := &ComplicatedStruct2{
		FloatArrayVal2: [][]float64{{1.0, 2.0}},
		StringValue1:   "line",
	}

	src := []*ComplicatedStruct1{
		{
			FloatValue1:           &FloatValue1,
			FloatValue2:           &FloatValue2,
			ComplicatedStruct2Ptr: ComplicatedStruct2Ptr,
			ComplicatedArray1Val: []*ComplicatedStruct3{
				{
					FloatValue1: FloatValue1,
					FloatValue2: FloatValue2,
					ComplicatedArray2Val: []*ComplicatedStruct4{
						{
							FloatValue1:           FloatValue1,
							FloatValue2:           FloatValue2,
							ComplicatedStruct2Ptr: ComplicatedStruct2Ptr,
							ComplicatedStruct5Ptr: &ComplicatedStruct5{
								StringValue5:   3.0,
								StringValue6:   4.0,
								FloatArrayVal1: []float64{6.0, 8.0},
								StringValue7:   "affine",
								StringValue1:   "polyline",
							},
							StringValue3: "GOD MODE",
							StringValue4: "like a boss",
						},
					},
					StringValue2: "Happy New Year!",
				},
			},
		},
	}

	dest := []*ComplicatedStruct1{}

	err := Copy(&dest, &src)
	if err != nil {
		t.Errorf("error copying data: %+v", err)
		return
	}

	if !reflect.DeepEqual(*dest[0], *src[0]) {
		t.Errorf("expected: %+v got: %+v", *src[0], *dest[0])
		return
	}
}

func TestInternalsExactValue(t *testing.T) {
	defer panicHandler(t)

	t1 := int64(666)
	t2 := &t1
	t3 := &t2
	t4 := &t3

	if exactValue(reflect.ValueOf(t4)).Kind() != reflect.Int64 {
		t.Errorf("exact value mismatch: expected %v got %+v", reflect.Int64, exactValue(reflect.ValueOf(t4)).Kind())
		return
	}
	if exactValue(reflect.ValueOf(t4)).Int() != t1 {
		t.Errorf("exact value mismatch: expected %d got %+v", t1, exactValue(reflect.ValueOf(t4)).Int())
		return
	}

	if exactValue(reflect.ValueOf(t2)).Kind() != reflect.Int64 {
		t.Errorf("exact value mismatch: expected %v got %+v", reflect.Int64, exactValue(reflect.ValueOf(t2)).Kind())
		return
	}
	if exactValue(reflect.ValueOf(t2)).Int() != t1 {
		t.Errorf("exact value mismatch: expected %d got %+v", t1, exactValue(reflect.ValueOf(t2)).Int())
		return
	}
}

func TestInternalsReductPointers(t *testing.T) {
	defer panicHandler(t)

	t1 := "la-la-la"
	t2 := &t1
	t3 := &t2
	t4 := &t3

	reducted := reductPointers(reflect.ValueOf(t4))

	if reducted.Kind() != reflect.Ptr {
		t.Errorf("wrong reducted kind")
		return
	}

	if reducted.Elem().String() != t1 {
		t.Errorf("wrong reducted value")
		return
	}
}

func TestInternalsDeepFields(t *testing.T) {
	defer panicHandler(t)

	dummy := TypeStruct3{}

	fields := deepFields(reflect.TypeOf(dummy))

	fieldsExpected := []string{"Field1", "Field2", "Field3", "Field4", "Field5"}

	if len(fields) != len(fieldsExpected) {
		t.Errorf("fields count mismatch")
		return
	}

	for index := range fields {
		if fields[index] != fieldsExpected[index] {
			t.Errorf("fields mismatch: expected %s got %s", fieldsExpected[index], fields[index])
			return
		}
	}

	dummy2 := &TypeStruct4{}
	fieldsExpected = []string{"Field2", "Field1"}

	fields = deepFields(reflect.TypeOf(dummy2))

	if len(fields) != len(fieldsExpected) {
		t.Errorf("fields count mismatch %v", fields)
		return
	}

	for index := range fields {
		if fields[index] != fieldsExpected[index] {
			t.Errorf("fields mismatch: expected %s got %s", fieldsExpected[index], fields[index])
			return
		}
	}
}

func TestInternalsFieldByName(t *testing.T) {
	defer panicHandler(t)

	dummyptr := &TypeStruct4{}

	dummy := &TypeStruct3{
		Field1: "razdvatri",
		Field2: "dldll",
		Field4: dummyptr,
		Field5: []*TypeStruct4{dummyptr},
	}

	dummy2 := &dummy
	dummy3 := &dummy2

	base := reflect.ValueOf(dummy3)

	f1 := fieldByName(base, "Field1")
	if f1.Interface().(string) != dummy.Field1 {
		t.Errorf("fields mismatch: expected %s got %s", dummy.Field1, f1.Interface())
		return
	}

	f2 := fieldByName(base, "Field2")
	if f2.String() != dummy.Field2 {
		t.Errorf("fields mismatch: expected %s got %s", dummy.Field2, f2.String())
		return
	}

	f4 := fieldByName(base, "Field4")
	if ((*TypeStruct4)(unsafe.Pointer(f4.Pointer()))) != dummy.Field4 {
		t.Errorf("fields mismatch: expected %+v got %+v", dummy.Field4, f4.Pointer())
		return
	}

	f5 := fieldByName(base, "Field5")
	if len(f5.Interface().([]*TypeStruct4)) != len(dummy.Field5) {
		t.Errorf("fields mismatch: expected %+v got %+v", dummy.Field5, f5.Interface().([]*TypeStruct4))
		return
	}

	if f5.Interface().([]*TypeStruct4)[0] != dummy.Field5[0] {
		t.Errorf("fields mismatch: expected %+v got %+v", dummy.Field5, f5.Interface().([]*TypeStruct4))
		return
	}
}

func TestInternalsMethodByName(t *testing.T) {
	defer panicHandler(t)

	dummyptr := &TypeStruct5{
		field1: "raz",
		Field2: "dva",
	}

	dummy := &dummyptr
	dummy2 := &dummy
	dummy3 := &dummy2

	base := reflect.ValueOf(dummy3)

	f1 := methodByName(base, "Field1")
	if !f1.IsValid() {
		t.Errorf("invalid method")
		return
	}

	expected := "tri"
	f1.Call([]reflect.Value{reflect.ValueOf(expected)})
	if dummyptr.field1 != expected {
		t.Errorf("wrong method results")
		return
	}
}

type DeepCopyPtrTest struct {
	Field1 **string
	Field2 interface{}
	Field3 *TypeStruct4
}

func TestInternalsDeepCopyPtr(t *testing.T) {
	defer panicHandler(t)

	dummyptr := &TypeStruct4{}

	strExpected := "razdvatri"
	strExpectedPtr := &strExpected

	dummy := &DeepCopyPtrTest{
		Field1: &strExpectedPtr,
		Field2: strExpectedPtr,
		Field3: dummyptr,
	}

	dummy2 := &dummy
	dummy3 := &dummy2

	dummyDest := &DeepCopyPtrTest{}
	dummyDest2 := &dummyDest
	dummyDest3 := &dummyDest2

	src := reflect.ValueOf(dummy3)
	dest := reflect.ValueOf(dummyDest3)

	d1 := fieldByName(dest, "Field1")
	s1 := fieldByName(src, "Field1")

	copied, err := tryDeepCopyPtr(d1, s1, nil)
	if !copied || err != nil {
		t.Errorf("error deep copy: copied %t err %v", copied, err)
		return
	}

	d2 := fieldByName(dest, "Field3")
	s2 := fieldByName(src, "Field3")

	copied, err = tryDeepCopyPtr(d2, s2, nil)
	if !copied || err != nil {
		t.Errorf("error deep copy: copied %t err %v", copied, err)
		return
	}

	if dummyDest.Field1 == nil || *dummyDest.Field1 == nil || **dummy.Field1 != **dummyDest.Field1 {
		t.Errorf("field mismatch: expected %v got %v", dummy.Field1, dummyDest.Field1)
	}

	if dummyDest.Field3 == nil || dummy.Field3.field1 != dummyDest.Field3.field1 {
		t.Errorf("field mismatch: expected %v got %v", dummy.Field3, dummyDest.Field3)
	}
}

func TestInternalsDeepCopyStruct(t *testing.T) {
	defer panicHandler(t)

	dummy := &TypeStruct1{
		Field3: TypeStruct2{
			Field1: 12,
			Field2: "tri",
		},
	}

	dummy2 := &dummy
	dummy3 := &dummy2

	dummyDest := &TypeStruct1{}
	dummyDest2 := &dummyDest
	dummyDest3 := &dummyDest2

	src := reflect.ValueOf(dummy3)
	dest := reflect.ValueOf(dummyDest3)

	d1 := fieldByName(dest, "Field3")
	s1 := fieldByName(src, "Field3")

	copied, err := tryDeepCopyStruct(d1, s1, nil)
	if !copied || err != nil {
		t.Errorf("error deep copy: copied %t err %v", copied, err)
	}
}

func panicHandler(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("test did panic: %v", r)
	}
}
