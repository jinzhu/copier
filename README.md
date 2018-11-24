# Copier

  I am a copier, I copy everything from one to another

[![wercker status](https://app.wercker.com/status/9d44ad2d4e6253929c8fb71359effc0b/s/master "wercker status")](https://app.wercker.com/project/byKey/9d44ad2d4e6253929c8fb71359effc0b)

## Features

* Copy from field to field with same name
* Copy from method to field with same name
* Copy from field to method with same name
* Copy from slice to slice
* Copy from struct to slice
* Extensible

## Usage

```go
package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

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

func main() {
	var (
		user      = User{Name: "Jinzhu", Age: 18, Role: "Admin"}
		users     = []User{{Name: "Jinzhu", Age: 18, Role: "Admin"}, {Name: "jinzhu 2", Age: 30, Role: "Dev"}}
		employee  = Employee{}
		employees = []Employee{}
	)

	copier.Copy(&employee, &user)

	fmt.Printf("%#v \n", employee)
	// Employee{
	//    Name: "Jinzhu",           // Copy from field
	//    Age: 18,                  // Copy from field
	//    DoubleAge: 36,            // Copy from method
	//    EmployeeId: 0,            // Ignored
	//    SuperRule: "Super Admin", // Copy to method
	// }

	// Copy struct to slice
	copier.Copy(&employees, &user)

	fmt.Printf("%#v \n", employees)
	// []Employee{
	//   {Name: "Jinzhu", Age: 18, DoubleAge: 36, EmployeId: 0, SuperRule: "Super Admin"}
	// }

	// Copy slice to slice
	employees = []Employee{}
	copier.Copy(&employees, &users)

	fmt.Printf("%#v \n", employees)
	// []Employee{
	//   {Name: "Jinzhu", Age: 18, DoubleAge: 36, EmployeId: 0, SuperRule: "Super Admin"},
	//   {Name: "jinzhu 2", Age: 30, DoubleAge: 60, EmployeId: 0, SuperRule: "Super Dev"},
	// }
}
```


## Usage with Extention CopierFunc
```go

package main

import (
	"encoding/json"
	"errors"
	"github.com/ariffebr/copier"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"reflect"
	"time"
)

type AccessToken struct {
	Id        string
	AuthToken string
	CreatedAt time.Time
	ExpiredAt time.Time
}

type ProtoAccessToken struct {
	Id        string
	AuthToken string
	CreatedAt *timestamp.Timestamp
	ExpiredAt *timestamp.Timestamp
}

// Convert from time.Time to Protobuf Timestamp
func fromTime2Timestamp(to, from reflect.Value) (err error) {
	log.Println(to.Addr().Type(), "->", from.Addr().Type())

	if _, ok := to.Addr().Interface().(*timestamp.Timestamp); ok {
		if fromTime, ok2 := from.Addr().Interface().(*time.Time); ok2 {
			var ts *timestamp.Timestamp
			ts, err = ptypes.TimestampProto(*fromTime)
			to.Set(reflect.Indirect(reflect.ValueOf(ts)))
		} else {
			err = errors.New("not from time.Time")
		}
	} else {
		err = errors.New("not to timestamp.Timestamp")
	}
	return err
}

// Convert from protobuf Timestamp to time.Time
func fromTimestamp2Time(to, from reflect.Value) (err error) {
	log.Println(to.Addr().Type(), "->", from.Addr().Type())

	if _, ok := to.Addr().Interface().(*time.Time); ok {
		if fromTimestamp, ok2 := from.Addr().Interface().(**timestamp.Timestamp); ok2 {
			var t time.Time
			t, err = ptypes.Timestamp(*fromTimestamp)

			to.Set(reflect.Indirect(reflect.ValueOf(t)))
		} else if fromTimestamp, ok2 := from.Addr().Interface().(*timestamp.Timestamp); ok2 {
			var t time.Time
			t, err = ptypes.Timestamp(fromTimestamp)
			to.Set(reflect.Indirect(reflect.ValueOf(t)))
		} else {
			err = errors.New("not from timestamp.Timestamp")
		}
	} else {
		err = errors.New("not to time.Time")
	}
	return err
}

func main() {

    // Register CopierFunc
	copier.RegisterCopyFunc(
		copier.CopierFunc{
			ToType:   reflect.TypeOf(timestamp.Timestamp{}),
			FromType: reflect.TypeOf(time.Time{}),
			CopyFunc: fromTime2Timestamp,
		},
		copier.CopierFunc{
			ToType:   reflect.TypeOf(time.Time{}),
			FromType: reflect.TypeOf(timestamp.Timestamp{}),
			CopyFunc: fromTimestamp2Time,
		},
	)

	accessToken := AccessToken{
		Id:        "ini token id - AccessToken",
		AuthToken: "ini auth tokennya - AccessToken",
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 24),
	}

	var pbAccessToken ProtoAccessToken
	err := copier.Copy(&pbAccessToken, &accessToken)

	if err != nil {
		log.Fatal("Failed to copy(1)", err)
	}

	resJson, _ := json.Marshal(pbAccessToken)
	log.Println("(1) copy[AccessToken -> ProtoAccessToken] result:\n\t", string(resJson))

	expired, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * 24))
	pbAccToken2 := ProtoAccessToken{
		Id:        "ini token_id - ProtoAccessToken",
		AuthToken: "ini auth token nya - ProtoAccessToken",
		CreatedAt: ptypes.TimestampNow(),
		ExpiredAt: expired,
	}

	var accessToken2 AccessToken

	err = copier.Copy(&accessToken2, &pbAccToken2)
	if err != nil {
		log.Println("Failed to copy(2)", err)
	}

	res2Json, _ := json.Marshal(accessToken2)
	log.Println("copy(2) [ProtoAccessToken -> AccessToken] result:\n\t", string(res2Json))

}


```

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the [MIT License](https://github.com/jinzhu/copier/blob/master/License).
