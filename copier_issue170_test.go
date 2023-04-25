package copier_test

import (
	"github.com/jinzhu/copier"
	"reflect"
	"testing"
)

type A struct {
	A int
}
type B struct {
	A int
	b int
}

var copied = B{A: 2387483274, b: 128387134}

func newOptWithConverter() copier.Option {
	return copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: A{},
				DstType: B{},
				Fn: func(from interface{}) (interface{}, error) {
					return copied, nil
				},
			},
		},
	}
}

func Test_Struct_With_Converter(t *testing.T) {
	aa := A{A: 11}
	bb := B{A: 10, b: 100}
	err := copier.CopyWithOption(&bb, &aa, newOptWithConverter())
	if err != nil || !reflect.DeepEqual(copied, bb) {
		t.Fatalf("Got %v, wanted %v", bb, copied)
	}
}

func Test_Map_With_Converter(t *testing.T) {
	aa := map[string]*A{
		"a": &A{A: 10},
	}

	bb := map[string]*B{
		"a": &B{A: 10, b: 100},
	}

	err := copier.CopyWithOption(&bb, &aa, newOptWithConverter())
	if err != nil {
		t.Fatalf("copy with converter failed: %v", err)
	}

	for _, v := range bb {
		wanted := &copied
		if !reflect.DeepEqual(v, wanted) {
			t.Fatalf("Got %v, wanted %v", v, wanted)
		}
	}
}

func Test_Slice_With_Converter(t *testing.T) {
	aa := []*A{
		&A{A: 10},
	}

	bb := []*B{
		&B{A: 10, b: 100},
	}

	err := copier.CopyWithOption(&bb, &aa, newOptWithConverter())

	if err != nil {
		t.Fatalf("copy slice error: %v", err)
	}

	wanted := copied
	for _, v := range bb {
		temp := v
		if !reflect.DeepEqual(*temp, wanted) {
			t.Fatalf("Got %v, wanted %v", *temp, wanted)
		}
	}
}

func Test_Slice_Embedded_With_Converter(t *testing.T) {
	aa := struct {
		A []*A
	}{
		A: []*A{&A{A: 10}},
	}

	bb := struct {
		A []*B
	}{
		A: []*B{&B{A: 10, b: 100}},
	}

	err := copier.CopyWithOption(&bb, &aa, newOptWithConverter())

	wanted := struct {
		A []*B
	}{
		A: []*B{&copied},
	}

	if err != nil || !reflect.DeepEqual(bb, wanted) {
		t.Fatalf("Got %v, wanted %v", bb, wanted)
	}
}
