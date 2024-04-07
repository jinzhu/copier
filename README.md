# Copier

I am a copier, I copy everything from one to another

[![test status](https://github.com/jinzhu/copier/workflows/tests/badge.svg?branch=master "test status")](https://github.com/jinzhu/copier/actions)

## Key Features

- Field-to-field and method-to-field copying based on matching names
- Support for copying data:
  - From slice to slice
  - From struct to slice
  - From map to map
- Field manipulation through tags:
  - Enforce field copying with `copier:"must"`
  - Override fields even when `IgnoreEmpty` is set with `copier:"override"`
  - Exclude fields from being copied with `copier:"-"`

## Getting Started

### Installation

To start using Copier, install Go and run go get:

```bash
go get -u github.com/jinzhu/copier
```

## Basic

Import Copier into your application to access its copying capabilities

```go
import "github.com/jinzhu/copier"
```

### Basic Copying

```go
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
	SuperRole string
}

func (employee *Employee) Role(role string) {
	employee.SuperRole = "Super " + role
}

func main() {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin"}
	employee := Employee{}

	copier.Copy(&employee, &user)
	fmt.Printf("%#v\n", employee)
	// Output: Employee{Name:"Jinzhu", Age:18, DoubleAge:36, SuperRole:"Super Admin"}
}
```

## Tag Usage Examples

### `copier:"-"` - Ignoring Fields

Fields tagged with `copier:"-"` are explicitly ignored by Copier during the copying process.

```go
type Source struct {
    Name   string
    Secret string // We do not want this to be copied.
}

type Target struct {
    Name   string
    Secret string `copier:"-"`
}

func main() {
    source := Source{Name: "John", Secret: "so_secret"}
    target := Target{}

    copier.Copy(&target, &source)
    fmt.Printf("Name: %s, Secret: '%s'\n", target.Name, target.Secret)
    // Output: Name: John, Secret: ''
}
```

### `copier:"must"` - Enforcing Field Copy

The `copier:"must"` tag forces a field to be copied, resulting in a panic or an error if the field cannot be copied.

```go
type MandatorySource struct {
	Identification int
}

type MandatoryTarget struct {
	ID int `copier:"must"` // This field must be copied, or it will panic/error.
}

func main() {
	source := MandatorySource{}
	target := MandatoryTarget{ID: 10}

	// This will result in a panic or an error since ID is a must field but is empty in source.
	if err := copier.Copy(&target, &source); err != nil {
		log.Fatal(err)
	}
}
```

### `copier:"must,nopanic"` - Enforcing Field Copy Without Panic

Similar to `copier:"must"`, but Copier returns an error instead of panicking if the field is not copied.

```go
type SafeSource struct {
	ID string
}

type SafeTarget struct {
	Code string `copier:"must,nopanic"` // Enforce copying without panic.
}

func main() {
	source := SafeSource{}
	target := SafeTarget{Code: "200"}

	if err := copier.Copy(&target, &source); err != nil {
		log.Fatalln("Error:", err)
	}
	// This will not panic, but will return an error due to missing mandatory field.
}
```

### `copier:"override"` - Overriding Fields with IgnoreEmpty

Fields tagged with `copier:"override"` are copied even if IgnoreEmpty is set to true in Copier options and works for nil values.

```go
type SourceWithNil struct {
    Details *string
}

type TargetOverride struct {
    Details *string `copier:"override"` // Even if source is nil, copy it.
}

func main() {
    details := "Important details"
    source := SourceWithNil{Details: nil}
    target := TargetOverride{Details: &details}

    copier.CopyWithOption(&target, &source, copier.Option{IgnoreEmpty: true})
    if target.Details == nil {
        fmt.Println("Details field was overridden to nil.")
    }
}
```

### Specifying Custom Field Names

Use field tags to specify a custom field name when the source and destination field names do not match.

```go
type SourceEmployee struct {
    Identifier int64
}

type TargetWorker struct {
    ID int64 `copier:"Identifier"` // Map Identifier from SourceEmployee to ID in TargetWorker
}

func main() {
    source := SourceEmployee{Identifier: 1001}
    target := TargetWorker{}

    copier.Copy(&target, &source)
    fmt.Printf("Worker ID: %d\n", target.ID)
    // Output: Worker ID: 1001
}
```

## Other examples

### Copy from Method to Field with Same Name

Illustrates copying from a method to a field and vice versa.

```go
// Assuming User and Employee structs defined earlier with method and field respectively.

func main() {
    user := User{Name: "Jinzhu", Age: 18}
    employee := Employee{}

    copier.Copy(&employee, &user)
    fmt.Printf("DoubleAge: %d\n", employee.DoubleAge)
    // Output: DoubleAge: 36, demonstrating method to field copying.
}
```

### Copy Struct to Slice

```go
func main() {
    user := User{Name: "Jinzhu", Age: 18, Role: "Admin"}
    var employees []Employee

    copier.Copy(&employees, &user)
    fmt.Printf("%#v\n", employees)
    // Output: []Employee{{Name: "Jinzhu", Age: 18, DoubleAge: 36, SuperRole: "Super Admin"}}
}
```

### Copy Slice to Slice

```go
func main() {
    users := []User{{Name: "Jinzhu", Age: 18, Role: "Admin"}, {Name: "jinzhu 2", Age: 30, Role: "Dev"}}
    var employees []Employee

    copier.Copy(&employees, &users)
    fmt.Printf("%#v\n", employees)
    // Output: []Employee{{Name: "Jinzhu", Age: 18, DoubleAge: 36, SuperRole: "Super Admin"}, {Name: "jinzhu 2", Age: 30, DoubleAge: 60, SuperRole: "Super Dev"}}
}
```

### Copy Map to Map

```go
func main() {
    map1 := map[int]int{3: 6, 4: 8}
    map2 := map[int32]int8{}

    copier.Copy(&map2, map1)
    fmt.Printf("%#v\n", map2)
    // Output: map[int32]int8{3:6, 4:8}
}
```

## Complex Data Copying: Nested Structures with Slices

This example demonstrates how Copier can be used to copy data involving complex, nested structures, including slices of structs, to showcase its ability to handle intricate data copying scenarios.

```go
package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type Address struct {
	City    string
	Country string
}

type Contact struct {
	Email  string
	Phones []string
}

type Employee struct {
	Name      string
	Age       int32
	Addresses []Address
	Contact   *Contact
}

type Manager struct {
	Name            string `copier:"must"`
	Age             int32  `copier:"must,nopanic"`
	ManagedCities   []string
	Contact         *Contact `copier:"override"`
	SecondaryEmails []string
}

func main() {
	employee := Employee{
		Name: "John Doe",
		Age:  30,
		Addresses: []Address{
			{City: "New York", Country: "USA"},
			{City: "San Francisco", Country: "USA"},
		},
		Contact: nil,
	}

	manager := Manager{
		ManagedCities: []string{"Los Angeles", "Boston"},
		Contact: &Contact{
			Email:  "john.doe@example.com",
			Phones: []string{"123-456-7890", "098-765-4321"},
		}, // since override is set this should be overridden with nil
		SecondaryEmails: []string{"secondary@example.com"},
	}

	copier.CopyWithOption(&manager, &employee, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	fmt.Printf("Manager: %#v\n", manager)
	// Output: Manager struct showcasing copied fields from Employee,
	// including overridden and deeply copied nested slices.
}
```

## Available tags

| Tag                 | Description                                                                                                       |
| ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `copier:"-"`        | Explicitly ignores the field during copying.                                                                      |
| `copier:"must"`     | Forces the field to be copied; Copier will panic or return an error if the field is not copied.                   |
| `copier:"nopanic"`  | Copier will return an error instead of panicking.                                                                 |
| `copier:"override"` | Forces the field to be copied even if `IgnoreEmpty` is set. Useful for overriding existing values with empty ones |
| `FieldName`         | Specifies a custom field name for copying when field names do not match between structs.                          |

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**jinzhu**

- <http://github.com/jinzhu>
- <wosmvp@gmail.com>
- <http://twitter.com/zhangjinzhu>

## License

Released under the [MIT License](https://github.com/jinzhu/copier/blob/master/License).
