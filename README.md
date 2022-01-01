# Copier

I am a copier, I copy everything from one to another

## Features

- Copy from field to field with same name
- Copy from method to field with same name
- Copy from field to method with same name
- Copy from slice to slice
- Copy from struct to slice
- Copy from map to map
- Enforce copying a field with a tag
- Ignore a field with a tag
- Deep Copy

## Usage

```go
package main

import (
	"fmt"
	"github.com/i-curve/copier"
)

type User struct {
	Name        string
	Role        string
	Age         int32
	EmployeCode int64 `copier:"EmployeNum"` // specify field name

	// Explicitly ignored in the destination struct.
	Salary   int
}

func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

// Tags in the destination Struct provide instructions to copier.Copy to ignore
// or enforce copying and to panic or return an error if a field was not copied.
type Employee struct {
	// Tell copier.Copy to panic if this field is not copied.
	Name      string `copier:"must"`

	// Tell copier.Copy to return an error if this field is not copied.
	Age       int32  `copier:"must,nopanic"`

	// Tell copier.Copy to explicitly ignore copying this field.
	Salary    int    `copier:"-"`

	DoubleAge int32
	EmployeId int64 `copier:"EmployeNum"` // specify field name
	SuperRole string
}

func (employee *Employee) Role(role string) {
	employee.SuperRole = "Super " + role
}

func main() {
	var (
		user      = User{Name: "i-curve", Age: 18, Role: "Admin", Salary: 200000}
		users     = []User{{Name: "i-curve", Age: 18, Role: "Admin", Salary: 100000}, {Name: "i-curve 2", Age: 30, Role: "Dev", Salary: 60000}}
		employee  = Employee{Salary: 150000}
		employees = []Employee{}
	)

	copier.Copy(&employee, &user)

	fmt.Printf("%#v \n", employee)
	// Employee{
	//    Name: "i-curve",           // Copy from field
	//    Age: 18,                  // Copy from field
	//    Salary:150000,            // Copying explicitly ignored
	//    DoubleAge: 36,            // Copy from method
	//    EmployeeId: 0,            // Ignored
	//    SuperRole: "Super Admin", // Copy to method
	// }

	// Copy struct to slice
	copier.Copy(&employees, &user)

	fmt.Printf("%#v \n", employees)
	// []Employee{
	//   {Name: "i-curve", Age: 18, Salary:0, DoubleAge: 36, EmployeId: 0, SuperRole: "Super Admin"}
	// }

	// Copy slice to slice
	employees = []Employee{}
	copier.Copy(&employees, &users)

	fmt.Printf("%#v \n", employees)
	// []Employee{
	//   {Name: "i-curve", Age: 18, Salary:0, DoubleAge: 36, EmployeId: 0, SuperRole: "Super Admin"},
	//   {Name: "i-curve 2", Age: 30, Salary:0, DoubleAge: 60, EmployeId: 0, SuperRole: "Super Dev"},
	// }

 	// Copy map to map
	map1 := map[int]int{3: 6, 4: 8}
	map2 := map[int32]int8{}
	copier.Copy(&map2, map1)

	fmt.Printf("%#v \n", map2)
	// map[int32]int8{3:6, 4:8}
}
```

### Copy with Option

```go
copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
```

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**i-curve**

- <http://github.com/i-curve>
- <wosmvp@gmail.com>
- <http://twitter.com/zhangi-curve>

## License

Released under the [MIT License](https://github.com/i-curve/copier/blob/master/License).
