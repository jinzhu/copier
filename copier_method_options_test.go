package copier_test

import (
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
)

type methodOptionSource struct {
	ignoredCalls int
}

func (s *methodOptionSource) Ignored() string {
	s.ignoredCalls++
	return "do not copy"
}

func (*methodOptionSource) Number() int        { return 0 }
func (*methodOptionSource) Empty() string      { return "" }
func (*methodOptionSource) Enabled() bool      { return false }
func (*methodOptionSource) Pointer() *int      { return nil }
func (*methodOptionSource) Text() string       { return "copied" }
func (*methodOptionSource) Overridden() string { return "" }

func TestCopyMethodSkipOptions(t *testing.T) {
	for _, tc := range []struct {
		name        string
		ignoreEmpty bool
		deepCopy    bool
	}{
		{"copy empty", false, false},
		{"ignore empty", true, false},
		{"deep copy empty", false, true},
		{"deep ignore empty", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			number := 42
			destination := struct {
				Ignored    string `copier:"-"`
				Number     int
				Empty      string
				Enabled    bool
				Pointer    *int
				Mapped     string `copier:"Text"`
				Overridden string `copier:"override"`
			}{"keep", 42, "keep", true, &number, "before", "clear"}
			want := destination
			want.Mapped = "copied"
			want.Overridden = ""
			if !tc.ignoreEmpty {
				want.Number = 0
				want.Empty = ""
				want.Enabled = false
				want.Pointer = nil
			}

			source := &methodOptionSource{}
			if err := copier.CopyWithOption(&destination, source, copier.Option{
				IgnoreEmpty: tc.ignoreEmpty,
				DeepCopy:    tc.deepCopy,
			}); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(destination, want) {
				t.Errorf("got %#v, want %#v", destination, want)
			}
			if source.ignoredCalls != 0 {
				t.Errorf("ignored method was called %d times", source.ignoredCalls)
			}
		})
	}
}
