package copier_test

import (
	"encoding/json"
	"github.com/ybzhanghx/copier"
	"testing"
)

func BenchmarkCopyStructTag(b *testing.B) {
	var fakeAge int32 = 12
	user := UserTag{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin",
		Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		copier.CopyByTag(&EmployeeTag{}, &user, "mson")
	}
}

func BenchmarkNamaCopyTag(b *testing.B) {
	var fakeAge int32 = 12
	user := UserTag{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin",
		Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		employee := &EmployeeTag{
			Name:      user.Name,
			Nickname:  &user.Nickname,
			Age:       int64(user.Age),
			FakeAge:   int(*user.FakeAge),
			DoubleAge: user.DoubleAge_int32(),
			Notes:     user.Notes,
		}
		employee.Role_string(user.Role)
	}
}

func BenchmarkJsonMarshalCopyTag(b *testing.B) {
	var fakeAge int32 = 12
	user := UserTag{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin",
		Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee EmployeeTag
		json.Unmarshal(data, &employee)

		employee.DoubleAge = user.DoubleAge_int32()
		employee.Role_string(user.Role)
	}
}
