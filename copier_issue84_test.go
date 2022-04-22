package copier_test

import (
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

type Embedded struct {
	Field1 string
	Field2 string
}

type Embedder struct {
	Embedded
	PtrField *string
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotWork struct {
	ID      string  `json:"id"`
	UserID  *string `json:"user_id"`
	Name    string  `json:"name"`
	Website *string `json:"website"`
	Timestamps
}

type Work struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	UserID  *string `json:"user_id"`
	Website *string `json:"website"`
	Timestamps
}

func TestIssue84(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		var embedder Embedder
		embedded := Embedded{
			Field1: "1",
			Field2: "2",
		}
		err := copier.Copy(&embedder, &embedded)
		if err != nil {
			t.Errorf("unable to copy: %s", err)
		}
		if embedder.Field1 != embedded.Field1 {
			t.Errorf("field1 value is %s instead of %s", embedder.Field1, embedded.Field1)
		}
		if embedder.Field2 != embedded.Field2 {
			t.Errorf("field2 value is %s instead of %s", embedder.Field2, embedded.Field2)
		}
	})
	t.Run("from issue", func(t *testing.T) {
		notWorkObj := NotWork{
			ID:      "123",
			Name:    "name",
			Website: nil,
			UserID:  nil,
			Timestamps: Timestamps{
				UpdatedAt: time.Now(),
			},
		}
		workObj := Work{
			ID:      "123",
			Name:    "name",
			Website: nil,
			UserID:  nil,
			Timestamps: Timestamps{
				UpdatedAt: time.Now(),
			},
		}

		destObj1 := Work{}
		destObj2 := NotWork{}

		copier.CopyWithOption(&destObj1, &workObj, copier.Option{IgnoreEmpty: true, DeepCopy: false})

		copier.CopyWithOption(&destObj2, &notWorkObj, copier.Option{IgnoreEmpty: true, DeepCopy: false})
	})
}
