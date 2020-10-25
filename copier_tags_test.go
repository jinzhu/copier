package copier

import "testing"

type EmployeeTags struct {
	Name    string `copier:"must"`
	DOB     string
	Address string
	ID      int `copier:"-"`
}

type User1 struct {
	Name    string
	DOB     string
	Address string
	ID      int
}

type User2 struct {
	DOB     string
	Address string
	ID      int
}

func TestCopyTagIgnore(t *testing.T) {
	employee := EmployeeTags{ID: 100}
	user := User1{Name: "Dexter Ledesma", DOB: "1 November, 1970", Address: "21 Jump Street", ID: 12345}
	Copy(&employee, user)
	if employee.ID == user.ID {
		t.Error("Was not expected to copy IDs")
	}
	if employee.ID != 100 {
		t.Error("Original ID was overwritten")
	}
}

func TestCopyTagMust(t *testing.T) {
	employee := &EmployeeTags{}
	user := &User2{DOB: "1 January 1970"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic.")
		}
	}()
	Copy(employee, user)
}
