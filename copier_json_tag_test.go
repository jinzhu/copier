package copier_test

import (
	"testing"

	"github.com/betterlmy/copier"
)

func TestCopyWithJSONTag(t *testing.T) {
	t.Run("basic json tag mapping", func(t *testing.T) {
		type Source struct {
			Name  string `json:"user_name"`
			Email string `json:"user_email"`
			Age   int    `json:"user_age"`
		}

		type Target struct {
			UserName  string `json:"user_name"`
			UserEmail string `json:"user_email"`
			UserAge   int    `json:"user_age"`
		}

		src := Source{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
		}

		var dst Target
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if dst.UserName != src.Name {
			t.Errorf("Expected UserName to be %s, got %s", src.Name, dst.UserName)
		}
		if dst.UserEmail != src.Email {
			t.Errorf("Expected UserEmail to be %s, got %s", src.Email, dst.UserEmail)
		}
		if dst.UserAge != src.Age {
			t.Errorf("Expected UserAge to be %d, got %d", src.Age, dst.UserAge)
		}
	})

	t.Run("json tag with omitempty", func(t *testing.T) {
		type Source struct {
			Name string `json:"full_name,omitempty"`
			Age  int    `json:"age,omitempty"`
		}

		type Target struct {
			FullName string `json:"full_name,omitempty"`
			Age      int    `json:"age,omitempty"`
		}

		src := Source{
			Name: "Jane Doe",
			Age:  25,
		}

		var dst Target
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if dst.FullName != src.Name {
			t.Errorf("Expected FullName to be %s, got %s", src.Name, dst.FullName)
		}
		if dst.Age != src.Age {
			t.Errorf("Expected Age to be %d, got %d", src.Age, dst.Age)
		}
	})

	t.Run("mixed json and copier tags", func(t *testing.T) {
		type Source struct {
			Name   string `json:"user_name"`
			Email  string `json:"email"`
			Status string `copier:"UserStatus"` // copier tag takes precedence
		}

		type Target struct {
			UserName   string `json:"user_name"`
			Email      string `json:"email"`
			UserStatus string
		}

		src := Source{
			Name:   "Bob Smith",
			Email:  "bob@example.com",
			Status: "active",
		}

		var dst Target
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if dst.UserName != src.Name {
			t.Errorf("Expected UserName to be %s, got %s", src.Name, dst.UserName)
		}
		if dst.Email != src.Email {
			t.Errorf("Expected Email to be %s, got %s", src.Email, dst.Email)
		}
		if dst.UserStatus != src.Status {
			t.Errorf("Expected UserStatus to be %s, got %s", src.Status, dst.UserStatus)
		}
	})

	t.Run("json tag ignored with dash", func(t *testing.T) {
		type Source struct {
			Name     string `json:"name"`
			Password string `json:"-"` // should be ignored by json tag
		}

		type Target struct {
			Name        string `json:"name"`
			SecretValue string `json:"-"` // different field name, won't be copied
		}

		src := Source{
			Name:     "Alice",
			Password: "secret123",
		}

		var dst Target
		dst.SecretValue = "original_secret"
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if dst.Name != src.Name {
			t.Errorf("Expected Name to be %s, got %s", src.Name, dst.Name)
		}
		// SecretValue should not be copied because field names don't match and json tag is "-"
		if dst.SecretValue != "original_secret" {
			t.Errorf("Expected SecretValue to remain as 'original_secret', got %s", dst.SecretValue)
		}
	})

	t.Run("without UseJSONTag option", func(t *testing.T) {
		type Source struct {
			Name string `json:"user_name"`
			Age  int    `json:"user_age"`
		}

		type Target struct {
			UserName string `json:"user_name"`
			UserAge  int    `json:"user_age"`
		}

		src := Source{
			Name: "Charlie",
			Age:  35,
		}

		var dst Target
		// Without UseJSONTag, should not match by json tags
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: false})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		// Fields should not be copied because field names don't match
		if dst.UserName != "" {
			t.Errorf("Expected UserName to be empty, got %s", dst.UserName)
		}
		if dst.UserAge != 0 {
			t.Errorf("Expected UserAge to be 0, got %d", dst.UserAge)
		}
	})

	t.Run("nested struct with json tags", func(t *testing.T) {
		type Address struct {
			Street string `json:"street_name"`
			City   string `json:"city_name"`
		}

		type SourceUser struct {
			Name    string  `json:"user_name"`
			Address Address `json:"address"`
		}

		type TargetUser struct {
			UserName string  `json:"user_name"`
			Address  Address `json:"address"`
		}

		src := SourceUser{
			Name: "David",
			Address: Address{
				Street: "123 Main St",
				City:   "New York",
			},
		}

		var dst TargetUser
		err := copier.CopyWithOption(&dst, &src, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if dst.UserName != src.Name {
			t.Errorf("Expected UserName to be %s, got %s", src.Name, dst.UserName)
		}
		if dst.Address.Street != src.Address.Street {
			t.Errorf("Expected Address.Street to be %s, got %s", src.Address.Street, dst.Address.Street)
		}
		if dst.Address.City != src.Address.City {
			t.Errorf("Expected Address.City to be %s, got %s", src.Address.City, dst.Address.City)
		}
	})

	t.Run("slice of structs with json tags", func(t *testing.T) {
		type Source struct {
			Name string `json:"user_name"`
			Age  int    `json:"user_age"`
		}

		type Target struct {
			UserName string `json:"user_name"`
			UserAge  int    `json:"user_age"`
		}

		sources := []Source{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
		}

		var targets []Target
		err := copier.CopyWithOption(&targets, &sources, copier.Option{UseJSONTag: true})
		if err != nil {
			t.Fatalf("Copy failed: %v", err)
		}

		if len(targets) != len(sources) {
			t.Fatalf("Expected %d targets, got %d", len(sources), len(targets))
		}

		for i := range sources {
			if targets[i].UserName != sources[i].Name {
				t.Errorf("Expected targets[%d].UserName to be %s, got %s", i, sources[i].Name, targets[i].UserName)
			}
			if targets[i].UserAge != sources[i].Age {
				t.Errorf("Expected targets[%d].UserAge to be %d, got %d", i, sources[i].Age, targets[i].UserAge)
			}
		}
	})
}
