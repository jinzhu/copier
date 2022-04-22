package copier_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

func TestCopyWithTypeConverters(t *testing.T) {
	type SrcStruct struct {
		Field1 time.Time
		Field2 *time.Time
		Field3 *time.Time
		Field4 string
	}

	type DestStruct struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 int
	}

	testTime := time.Date(2021, 3, 5, 1, 30, 0, 123000000, time.UTC)

	src := SrcStruct{
		Field1: testTime,
		Field2: &testTime,
		Field3: nil,
		Field4: "9000",
	}

	var dst DestStruct

	err := copier.CopyWithOption(&dst, &src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)

					if !ok {
						return nil, errors.New("src type not matching")
					}

					return s.Format(time.RFC3339), nil
				},
			},
			{
				SrcType: copier.String,
				DstType: copier.Int,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)

					if !ok {
						return nil, errors.New("src type not matching")
					}

					return strconv.Atoi(s)
				},
			},
		},
	})

	if err != nil {
		t.Fatalf(`Should be able to copy from src to dst object. %v`, err)
		return
	}

	dateStr := "2021-03-05T01:30:00Z"

	if dst.Field1 != dateStr {
		t.Fatalf("got %q, wanted %q", dst.Field1, dateStr)
	}

	if dst.Field2 != dateStr {
		t.Fatalf("got %q, wanted %q", dst.Field2, dateStr)
	}

	if dst.Field3 != "" {
		t.Fatalf("got %q, wanted %q", dst.Field3, "")
	}

	if dst.Field4 != 9000 {
		t.Fatalf("got %q, wanted %q", dst.Field4, 9000)
	}
}

func TestCopyWithConverterAndAnnotation(t *testing.T) {
	type SrcStruct struct {
		Field1 string
	}

	type DestStruct struct {
		Field1 string
		Field2 string `copier:"Field1"`
	}

	src := SrcStruct{
		Field1: "test",
	}

	var dst DestStruct

	err := copier.CopyWithOption(&dst, &src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: copier.String,
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)

					if !ok {
						return nil, errors.New("src type not matching")
					}

					return s + "2", nil
				},
			},
		},
	})

	if err != nil {
		t.Fatalf(`Should be able to copy from src to dst object. %v`, err)
		return
	}

	if dst.Field2 != "test2" {
		t.Fatalf("got %q, wanted %q", dst.Field2, "test2")
	}
}

func TestCopyWithConverterStrToStrPointer(t *testing.T) {
	type SrcStruct struct {
		Field1 string
	}

	type DestStruct struct {
		Field1 *string
	}

	src := SrcStruct{
		Field1: "",
	}

	var dst DestStruct

	ptrStrType := ""

	err := copier.CopyWithOption(&dst, &src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: copier.String,
				DstType: &ptrStrType,
				Fn: func(src interface{}) (interface{}, error) {
					s, _ := src.(string)

					// return nil on empty string
					if s == "" {
						return nil, nil
					}

					return &s, nil
				},
			},
		},
	})

	if err != nil {
		t.Fatalf(`Should be able to copy from src to dst object. %v`, err)
		return
	}

	if dst.Field1 != nil {
		t.Fatalf("got %q, wanted nil", *dst.Field1)
	}
}
