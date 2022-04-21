# Copier[^1]

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
- Support from struct to map
- Support int64 timestamp and go time each copy

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

Copy struct to map

```go
	// Copy struct to map
	map3 := map[string]interface{}
	copier.Copy(&map3, &user)
	fmt.Println(map3)
	// map[age:18 name:i-curve role:Admin salary:200000]

	// Copy with options
	map4 := make(map[string]interface{})
	copier.CopyWithOption(&map4, &user, copier.Option{
		UpperCase:   false, //默认false, 是否设置拷贝到map的键转换为小写, 如果为true则不进行小写转化
		IgnoreEmpty: true, // 是否忽略空值, Copy中默认是true,CopyWithOption中默认是false,需要显示指定忽略空值
		IgnoreField: []string{"Age"}}) // 会忽略的字段范围, 如果UpperCase 为ture,这里也需要相应大写
	fmt.Println(map4)
	// map[Name:i-curve Role:Admin Salary:200000]
```

Copy int64 timestamp and go time

```go
	// 单值copier需要显示指定TimeFormat参数
	now := time.Now()
	var c int64
	copier.CopyWithOption(&c, &now, copier.Option{TimeFormat: "unixmill"})
	fmt.Println(now, now.UnixMilli(), c) // 毫秒
	// 2022-04-21 09:50:36.529860392 +0800 CST m=+0.000259408 1650505836529 1650505836529
	copier.CopyWithOption(&c, &now, copier.Option{TimeFormat: "unix"}) // 秒
	// 2022-04-21 09:51:31.793835808 +0800 CST m=+0.000226917 1650505891793 1650505891
	var timstreap int64 = 1650505969162
	b := time.Time{}
	copier.CopyWithOption(&b, &timstreap, copier.Option{TimeFormat: "unixmill"})
	fmt.Println(b)
	// 2022-04-21 09:52:49.162 +0800 CST

	// Struct: 详情请看copier_curve_test.go
	type TimInt struct {
		Time1 int64
		Time2 int64
		Time3 int64
		Time4 int64
		Time5 *int64
		Time6 *int64
		Time7 int64
	}
	type TimTim struct {
		Time1 time.Time  `copier:"time_format:unix"`
		Time2 time.Time  `copier:"-,time_format:unix"`
		Time3 *time.Time `copier:"time_format:unix"`
		Time4 *time.Time `copier:"time_format:unixmill"`
		Time5 time.Time  `copier:"time_format:unix"`
		Time6 time.Time  `copier:"time_format:unix"`
		Time7 time.Time
	}
	now := time.Now()
	var tim = TimTim{
		Time1: now,
		Time2: now,
		Time3: &now,
		Time4: &now,
	}
	var num TimInt
	copier.Copy(&num, &tim)
	fmt.Printf("%+v\n", num)
	// {Time1:1650506317 Time2:0 Time3:1650506317 Time4:1650506317495 Time5:<nil> Time6:<nil> Time7:0}
```

### Copy with Option

```go
	copier.CopyWithOption(&to, &from,
		copier.Option{
			IgnoreEmpty: true, // 是否忽略空值, copier中为true
			DeepCopy: true,
			UpperCase:  false, // struct to map拷贝时, 设定键名是大驼峰还是下划线小写
			IgnoreField: []string{""}, // 忽略拷贝的字段集
			TimeFormat: "unix", // 时间和int64拷贝时, 时间格式.只支持: unix(秒), unixmill(毫秒)
			TagFlag: "copier", // tag的名, 默认为copier, 可以设置为其他名字如:json
			TagDelimiter: ",", // tao内容的分隔符, 默认为","
		})
```

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**i-curve**

- <http://github.com/i-curve>
- <i-curve@qq.com>

## License

Released under the [MIT License](https://github.com/i-curve/copier/blob/master/License).

[^1]: fork 自 [https://github.com/jinzhu/copier](https://github.com/jinzhu/copier),由于使用习惯添加了部分功能但并没有被成功合并进去,两个项目使用时已有些许差异
