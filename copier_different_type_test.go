package copier_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

type TypeStruct1 struct {
	Field1 string
	Field2 string
	Field3 TypeStruct2
	Field4 *TypeStruct2
	Field5 []*TypeStruct2
	Field6 []TypeStruct2
	Field7 []*TypeStruct2
	Field8 []TypeStruct2
	Field9 []string
}

type TypeStruct2 struct {
	Field1 int
	Field2 string
	Field3 []TypeStruct2
	Field4 *TypeStruct2
	Field5 *TypeStruct2
	Field9 string
}

type TypeStruct3 struct {
	Field1 interface{}
	Field2 string
	Field3 TypeStruct4
	Field4 *TypeStruct4
	Field5 []*TypeStruct4
	Field6 []*TypeStruct4
	Field7 []TypeStruct4
	Field8 []TypeStruct4
}

type TypeStruct4 struct {
	field1 int
	Field2 string
}

func (t *TypeStruct4) Field1(i int) {
	t.field1 = i
}

type TypeBaseStruct5 struct {
	A bool
	B byte
	C float64
	D int16
	E int32
	F int64
	G time.Time
	H string
}

type TypeSqlNullStruct6 struct {
	A sql.NullBool    `json:"a"`
	B sql.NullByte    `json:"b"`
	C sql.NullFloat64 `json:"c"`
	D sql.NullInt16   `json:"d"`
	E sql.NullInt32   `json:"e"`
	F sql.NullInt64   `json:"f"`
	G sql.NullTime    `json:"g"`
	H sql.NullString  `json:"h"`
}

func TestCopyDifferentFieldType(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts2 := &TypeStruct2{}

	copier.Copy(ts2, ts)

	if ts2.Field2 != ts.Field2 || ts2.Field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts2")
	}
}

func TestCopyDifferentTypeMethod(t *testing.T) {
	ts := &TypeStruct1{
		Field1: "str1",
		Field2: "str2",
	}
	ts4 := &TypeStruct4{}

	copier.Copy(ts4, ts)

	if ts4.Field2 != ts.Field2 || ts4.field1 != 0 {
		t.Errorf("Should be able to copy from ts to ts4")
	}
}

func TestAssignableType(t *testing.T) {
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
		Field6: []TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
		Field7: []*TypeStruct2{
			{
				Field1: 666,
				Field2: "str2",
			},
		},
	}

	ts3 := &TypeStruct3{}

	copier.CopyWithOption(&ts3, &ts, copier.Option{CaseSensitive: true})

	if v, ok := ts3.Field1.(string); !ok {
		t.Error("Assign to interface{} type did not succeed")
	} else if v != "str1" {
		t.Error("String haven't been copied correctly")
	}

	if ts3.Field2 != ts.Field2 {
		t.Errorf("Field2 should be copied")
	}

	checkType2WithType4(ts.Field3, ts3.Field3, t, "Field3")
	checkType2WithType4(*ts.Field4, *ts3.Field4, t, "Field4")

	if len(ts3.Field5) != len(ts.Field5) {
		t.Fatalf("fields not equal, got %v, expects: %v", len(ts3.Field5), len(ts.Field5))
	}

	for idx, f := range ts.Field5 {
		checkType2WithType4(*f, *(ts3.Field5[idx]), t, "Field5")
	}

	for idx, f := range ts.Field6 {
		checkType2WithType4(f, *(ts3.Field6[idx]), t, "Field6")
	}

	for idx, f := range ts.Field7 {
		checkType2WithType4(*f, ts3.Field7[idx], t, "Field7")
	}

	for idx, f := range ts.Field8 {
		checkType2WithType4(f, ts3.Field8[idx], t, "Field8")
	}
}

func checkType2WithType4(t2 TypeStruct2, t4 TypeStruct4, t *testing.T, testCase string) {
	if t2.Field1 != t4.field1 || t2.Field2 != t4.Field2 {
		t.Errorf("%v: type struct 4 and type struct 2 is not equal", testCase)
	}
}

func TestCopyFromBaseToSqlNullWithOptionDeepCopy(t *testing.T) {
	a := TypeBaseStruct5{
		A: true,
		B: byte(2),
		C: 5.5,
		D: 1,
		E: 2,
		F: 3,
		G: time.Now(),
		H: "deep",
	}
	b := TypeSqlNullStruct6{}

	err := copier.CopyWithOption(&b, a, copier.Option{DeepCopy: true})
	// 检查是否有错误
	if err != nil {
		t.Errorf("CopyStructWithOption() error = %v", err)
		return
	}
	// 检查 b 结构体的字段是否符合预期
	if !b.A.Valid || b.A.Bool != true {
		t.Errorf("b.A = %v, want %v", b.A, true)
	}
	if !b.B.Valid || b.B.Byte != byte(2) {
		t.Errorf("b.B = %v, want %v", b.B, byte(2))
	}
	if !b.C.Valid || b.C.Float64 != 5.5 {
		t.Errorf("b.C = %v, want %v", b.C, 5.5)
	}
	if !b.D.Valid || b.D.Int16 != 1 {
		t.Errorf("b.D = %v, want %v", b.D, 1)
	}
	if !b.E.Valid || b.E.Int32 != 2 {
		t.Errorf("b.E = %v, want %v", b.E, 2)
	}
	if !b.F.Valid || b.F.Int64 != 3 {
		t.Errorf("b.F = %v, want %v", b.F, 3)
	}
	if !b.G.Valid || b.G.Time != a.G {
		t.Errorf("b.G = %v, want %v", b.G, a.G)
	}
	if !b.H.Valid || b.H.String != "deep" {
		t.Errorf("b.H = %v, want %v", b.H, "deep")
	}
}
