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

func BenchmarkCopyStructFields(b *testing.B) {
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
			NickName:  &user.Nickname,
			Age:       int64(user.Age),
			FakeAge:   int(*user.FakeAge),
			DoubleAge: user.DoubleAge(),
		}

		employee.Notes = make([]*string, len(user.Notes))
		for idx, note := range user.Notes {
			tmp := note
			employee.Notes[idx] = &tmp
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
