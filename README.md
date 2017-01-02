# Copier

  I am a copier, I copy everything from one to another

## Features

* Copy struct's field to field if its name match
* Copy from method to field if its name match
* Copy from field to method if its name match
* Copy from slice to slice
* Copy from struct to slice

## Usage

```go
import "github.com/jinzhu/copier"

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

var (
  user     = User{Name: "Jinzhu", Age: 18, Role: "Admin"}
  employee = Employee{}
)

coper.Copy(&employee, &user)
// employee => Employee{ Name: "Jinzhu",           // Copy from field
//                       Age: 18,                  // Copy from field
//                       DoubleAge: 36,            // Copy from method
//                       EmployeeId: 0,            // Just ignored
//                       SuperRule: "Super Admin", // Copy to method
//                      }

// Copy struct to slice
var (
  user      = User{Name: "hello", Age: 18, Role: "User"}
  employees = []Employee{}
)

coper.Copy(&employees, &user)
// employees => [{hello 18 0 36 Super User}]


// Copy slice to slice
var (
  users     = []User{{Name: "Jinzhu", Age: 18, Role: "Admin"}, {Name: "jinzhu 2", Age: 30, Role: "Dev"}}
  employees = []Employee{}
)

coper.Copy(&employees, &users)
// employees => [{hello 18 0 36 Super User} {Jinzhu 18 0 36 Super Admin} {jinzhu 2 30 0 60 Super Dev}]
```

# Supporting the project

[![http://patreon.com/jinzhu](http://patreon_public_assets.s3.amazonaws.com/sized/becomeAPatronBanner.png)](http://patreon.com/jinzhu)

# Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the [MIT License](https://github.com/jinzhu/copier/blob/master/License).
