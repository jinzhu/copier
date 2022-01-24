package domain

import "testing"

func TestParseMaterialGroupType(t *testing.T) {
	mgType, err := ParseMaterialGroupType("abc")
	if err != nil {

	}
	t.Logf("mgType: %v", mgType)
}
