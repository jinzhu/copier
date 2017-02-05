package copier_test

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/copier"
)

func BenchmarkCopyStruct(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		copier.Copy(&Employee{}, &user)
	}
}

func BenchmarkNamaCopy(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			Name:      user.Name,
			Nickname:  &user.Nickname,
			Age:       int64(user.Age),
			FakeAge:   int(*user.FakeAge),
			DoubleAge: user.DoubleAge(),
			Notes:     user.Notes,
		}
		employee.Role(user.Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)

		employee.DoubleAge = user.DoubleAge()
		employee.Role(user.Role)
	}
}
