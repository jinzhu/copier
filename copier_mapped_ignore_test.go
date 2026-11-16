package copier_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
)

type mappedIgnoreSource struct {
	ID   int
	Id   int `copier:"ID"`
	Name string
}

type mappedIgnoreDestination struct {
	ID   int `copier:"-"`
	Name string
}

type MappedIgnoreEmbedded struct {
	ID int `copier:"-"`
}

func TestCopyMappedDestinationFlags(t *testing.T) {
	for _, deepCopy := range []bool{false, true} {
		for _, test := range []struct {
			name        string
			source      interface{}
			destination interface{}
			want        interface{}
			option      copier.Option
		}{
			{
				name:        "source tag cannot bypass ignored destination",
				source:      mappedIgnoreSource{ID: 111, Id: 12345, Name: "copied"},
				destination: &mappedIgnoreDestination{ID: 100},
				want:        &mappedIgnoreDestination{ID: 100, Name: "copied"},
			},
			{
				name:        "explicit mapping cannot bypass ignored destination",
				source:      struct{ Identifier int }{12345},
				destination: &mappedIgnoreDestination{ID: 100},
				want:        &mappedIgnoreDestination{ID: 100},
				option: copier.Option{FieldNameMapping: []copier.FieldNameMapping{{
					SrcType: struct{ Identifier int }{}, DstType: mappedIgnoreDestination{},
					Mapping: map[string]string{"Identifier": "ID"},
				}}},
			},
			{
				name:        "case insensitive mapping cannot bypass ignored destination",
				source:      struct{ Identifier int }{12345},
				destination: &mappedIgnoreDestination{ID: 100},
				want:        &mappedIgnoreDestination{ID: 100},
				option: copier.Option{FieldNameMapping: []copier.FieldNameMapping{{
					SrcType: struct{ Identifier int }{}, DstType: mappedIgnoreDestination{},
					Mapping: map[string]string{"Identifier": "id"},
				}}},
			},
			{
				name:   "ignored promoted field does not initialize its parent",
				source: mappedIgnoreSource{Id: 12345, Name: "copied"},
				destination: &struct {
					*MappedIgnoreEmbedded
					Name string
				}{},
				want: &struct {
					*MappedIgnoreEmbedded
					Name string
				}{Name: "copied"},
			},
			{
				name: "ignored source-named destination does not hide a different target",
				source: struct {
					ID int `copier:"Number"`
				}{12345},
				destination: &struct {
					ID     int `copier:"-"`
					Number int
				}{ID: 100},
				want: &struct {
					ID     int `copier:"-"`
					Number int
				}{ID: 100, Number: 12345},
			},
			{
				name: "required mapped target is marked copied",
				source: struct {
					Identifier int `copier:"Number"`
				}{12345},
				destination: &struct {
					Number int `copier:"must,nopanic"`
				}{},
				want: &struct {
					Number int `copier:"must,nopanic"`
				}{Number: 12345},
			},
			{
				name: "exact case still takes precedence",
				source: struct {
					Identifier int `copier:"Id"`
				}{12345},
				destination: &struct {
					ID int `copier:"-"`
					Id int
				}{ID: 100},
				want: &struct {
					ID int `copier:"-"`
					Id int
				}{ID: 100, Id: 12345},
			},
		} {
			t.Run(fmt.Sprintf("%s/deep=%t", test.name, deepCopy), func(t *testing.T) {
				test.option.DeepCopy = deepCopy
				if err := copier.CopyWithOption(test.destination, test.source, test.option); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(test.destination, test.want) {
					t.Errorf("got %#v, want %#v", test.destination, test.want)
				}
			})
		}
	}
}
