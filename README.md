# Copier

  I am a copier, I copy everything from a struct to another struct

## Features

* Copy field to filed if name exactly match
* Copy from method to field if method name and field name exactly match
* Copy from field to method if field name and method name exactly match
* Copy slice to slice

## Usage

```go
import . "github.com/jinzhu/copier"

type User struct {
	Name string
	Role string
	Age  int32
}

func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Age       int32
	DoubleAge int32
	EmployeId int64
	SuperRule string
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

user := User{Name: "Jinzhu", Age: 18, Role: "Admin"}
employee := Employee{}

Copy(&employee, &user)

// employee => Employee{ Name: "Jinzhu",           // Copy from field
//                       Age: 18,                  // Copy from field
//                       DoubleAge: 36,            // Copy from method
//                       EmployeeId: 0,            // Just ignored
//                       SuperRule: "Super Admin", // Copy to method
//                      }

// Copy a slice? that's it
users := []User{{Name: "Jinzhu", Age: 18, Role: "Admin"}, {Name: "jinzhu 2", Age: 30, Role: "Dev"}}
employees := []Employee{}

Copy(&employees, &users)

```
