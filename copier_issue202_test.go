package copier_test

import (
	"testing"

	"github.com/jinzhu/copier"
)

type issue202Common struct {
	Common1 int
	Common2 string
}

type issue202Source struct {
	Name        string
	Description string
	issue202Common
}

type issue202Dest struct {
	Name        string
	Description string
	Common1     int
	Common2     string
}

func TestCopyFromUnexportedEmbeddedStruct(t *testing.T) {
	src := issue202Source{
		Name:           "name",
		Description:    "description",
		issue202Common: issue202Common{Common1: 2, Common2: "common2"},
	}
	dest := issue202Dest{}
	if err := copier.Copy(&dest, src); err != nil {
		t.Fatal(err)
	}
	want := issue202Dest{Name: "name", Description: "description", Common1: 2, Common2: "common2"}
	if dest != want {
		t.Errorf("got %+v, want %+v", dest, want)
	}
}
