package copier_test

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/copier"
)

func BenchmarkCopyStruct(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		copier.Copy(&Employee{}, &user)
	}
}

func BenchmarkNamaCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			Name:      user.Name,
			Age:       int64(user.Age),
			DoubleAge: user.DoubleAge(),
			Notes:     user.Notes,
		}
		employee.Role(user.Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)
		employee.DoubleAge = user.DoubleAge()
		employee.Role(user.Role)
	}
}
