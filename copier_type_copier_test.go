package copier_test

import (
	"errors"
	"github.com/jinzhu/copier"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type intToStringCopier struct{}

func (s intToStringCopier) Pairs() []copier.TypePair {
	return []copier.TypePair{
		{
			SrcType: reflect.TypeOf(0),
			DstType: reflect.TypeOf(""),
		},
	}
}

func (s intToStringCopier) Copy(dst, src reflect.Value) error {
	val, ok := src.Interface().(int)
	if !ok {
		return errors.New("type not match")
	}
	str := strconv.Itoa(val)
	dst.Set(reflect.ValueOf(str))
	return nil
}

type timeToStringCopier struct{}

func (t timeToStringCopier) Pairs() []copier.TypePair {
	return []copier.TypePair{
		{
			SrcType: reflect.TypeOf(time.Time{}),
			DstType: reflect.TypeOf(""),
		},
		{
			SrcType: reflect.TypeOf(&time.Time{}),
			DstType: reflect.TypeOf(""),
		},
	}
}

func (t timeToStringCopier) Copy(dst, src reflect.Value) error {
	const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"
	if src.Kind() == reflect.Ptr && src.IsNil() {
		if dst.Kind() == reflect.Ptr {
			dst.Set(reflect.Zero(reflect.TypeOf("")))
		}
		return nil
	}

	var val string
	if src.Kind() == reflect.Ptr {
		s, ok := src.Interface().(*time.Time)
		if !ok {
			return errors.New("type not match")
		}
		val = s.Format(timeFormat)
	} else {
		s, ok := src.Interface().(time.Time)
		if !ok {
			return errors.New("type not match")
		}
		val = s.Format(timeFormat)
	}
	dst.Set(reflect.ValueOf(val))
	return nil
}

func TestCopy_Register(t *testing.T) {
	type SrcStruct1 struct {
		Field1 int
		Field2 time.Time
		Field3 *time.Time
	}

	type DestStruct1 struct {
		Field1 string
		Field2 string
		Field3 string
	}

	type SrcStruct2 struct {
		Field1 SrcStruct1
		Field2 *SrcStruct1
		Field3 []SrcStruct1
		Field4 []*SrcStruct1
		Field5 map[int]SrcStruct1
		Field6 map[int]*SrcStruct1
	}

	type DestStruct2 struct {
		Field1 DestStruct1
		Field2 *DestStruct1
		Field3 []DestStruct1
		Field4 []*DestStruct1
		Field5 map[int]DestStruct1
		Field6 map[int]*DestStruct1
	}

	t.Run("copy different types", func(t *testing.T) {
		c := copier.NewCopier()
		c.Register(&intToStringCopier{})
		c.Register(&timeToStringCopier{})

		testTime := time.Date(2021, 3, 5, 1, 30, 0, 123000000, time.UTC)
		src := SrcStruct1{
			Field1: 100,
			Field2: testTime,
			Field3: &testTime,
		}
		var dst DestStruct1

		err := c.Copy(&dst, src)
		if err != nil {
			t.Error("copy fail")
			return
		}
		if dst.Field1 != "100" {
			t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
		}
		if dst.Field2 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
		}
		if dst.Field3 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
		}
	})

	t.Run("copy different types in map, slice, struct", func(t *testing.T) {
		c := copier.NewCopier()
		c.Register(&intToStringCopier{})
		c.Register(&timeToStringCopier{})

		testTime := time.Date(2021, 3, 5, 1, 30, 0, 123000000, time.UTC)
		src := SrcStruct2{
			Field1: SrcStruct1{
				Field1: 100,
				Field2: testTime,
				Field3: &testTime,
			},
			Field2: &SrcStruct1{
				Field1: 100,
				Field2: testTime,
				Field3: &testTime,
			},
			Field3: []SrcStruct1{
				{
					Field1: 100,
					Field2: testTime,
					Field3: &testTime,
				},
			},
			Field4: []*SrcStruct1{
				{
					Field1: 100,
					Field2: testTime,
					Field3: &testTime,
				},
			},
			Field5: map[int]SrcStruct1{
				1: {
					Field1: 100,
					Field2: testTime,
					Field3: &testTime,
				},
			},
			Field6: map[int]*SrcStruct1{
				1: {
					Field1: 100,
					Field2: testTime,
					Field3: &testTime,
				},
			},
		}
		var dst DestStruct2

		err := c.Copy(&dst, src)
		if err != nil {
			t.Error("copy fail")
			return
		}

		if dst.Field1.Field1 != "100" {
			t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
		}
		if dst.Field1.Field2 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
		}
		if dst.Field1.Field3 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
		}

		if dst.Field2.Field1 != "100" {
			t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
		}
		if dst.Field2.Field2 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
		}
		if dst.Field2.Field3 != "2021-03-05T01:30:00.123Z" {
			t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
		}

		for _, f := range dst.Field3 {
			if f.Field1 != "100" {
				t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
			}
			if f.Field2 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
			}
			if f.Field3 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
			}
		}

		for _, f := range dst.Field4 {
			if f.Field1 != "100" {
				t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
			}
			if f.Field2 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
			}
			if f.Field3 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
			}
		}

		for _, f := range dst.Field5 {
			if f.Field1 != "100" {
				t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
			}
			if f.Field2 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
			}
			if f.Field3 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
			}
		}

		for _, f := range dst.Field6 {
			if f.Field1 != "100" {
				t.Errorf("TestCopy_RegisterField: copy Field1 failed [%v]", dst.Field1)
			}
			if f.Field2 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field2 failed [%v]", dst.Field2)
			}
			if f.Field3 != "2021-03-05T01:30:00.123Z" {
				t.Errorf("TestCopy_RegisterField: copy Field3 failed [%v]", dst.Field3)
			}
		}

	})
}
