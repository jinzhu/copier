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
