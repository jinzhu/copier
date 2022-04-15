package copier_test

import (
	"testing"

	"github.com/i-curve/copier"
)

func TestGenernal(t *testing.T) {
	type stru struct {
		A int
		B string
	}
	var a = stru{
		A: 1,
		B: "2",
	}
	var b = stru{}
	copier.Copy(&a, &b)
	if a == b {
		t.Error("Copy copier the default value")
	}
	copier.CopyWithOption(&a, &b, copier.Option{IgnoreEmpty: false})
	if a != b {
		t.Error("CopyWithOption not copier the default value")
	}
	c := []int{1, 2, 3}
	var d []int
	copier.Copy(&c, &d)
	if len(c) == 0 {
		t.Error("Copy copier the default value")
	}
	copier.CopyWithOption(&c, &d, copier.Option{IgnoreEmpty: false})
	if len(c) != 0 {
		t.Error("CopyWithOption not copier the default value")
	}
}

func TestStringToLower(t *testing.T) {
	a := "NameIDc"
	if copier.StringToLower(a) != "name_idc" {
		t.Error("error of StringToLower")
	}
	a = "ClueIpr"
	if copier.StringToLower(a) != "clue_ipr" {
		t.Error("error of StringToLower")
	}
	a = "NameID"
	if copier.StringToLower(a) != "name_id" {
		t.Error("error of StringToLower")
	}
	a = "Address"
	if copier.StringToLower(a) != "address" {
		t.Error("error of StringToLower")
	}
	a = "addr_idess"
	if copier.StringToLower(a) != "addr_idess" {
		t.Error("error of StringToLower")
	}
}

func TestCopyFromStructToMap(t *testing.T) {
	var a = struct {
		NameID   int
		Username string
		Password []int `copier:"-"`
		Address  map[string]string
		Data     struct {
			Weight float64
			High   int
		}
	}{
		NameID:   1,
		Username: "i-curve",
		Password: []int{1, 2, 3},
		Address:  map[string]string{"home": "beijing"},
	}
	var b = map[string]interface{}{}
	if err := copier.CopyWithOption(&b, &a, copier.Option{IgnoreEmpty: false, IgnoreField: []string{"username"}}); err != nil {
		t.Errorf("Copy error from struct to map: %s", err)
	}
	if _, ok := b["username"]; ok {
		t.Error("Copy error from struct to map: field username should not be copied")
	}
	if _, ok := b["password"]; ok {
		t.Error("Copy error from struct to map: field password should not be copied")
	}
	if _, ok := b["address"]; !ok {
		t.Error("Copy error from struct to map: field address should be copied")
	}
}
