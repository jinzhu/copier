package copier_test

import (
	"fmt"
	"testing"
	"time"

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

func TestCopyTimeAndInt64(t *testing.T) {
	now := time.Now()
	var b time.Time
	copier.Copy(&b, &now)
	if b != now {
		t.Error("error of Copy time to time")
	}
	var c int64
	copier.CopyWithOption(&c, &now, copier.Option{TimeFormat: "unixmill"})
	if now.UnixMilli() != c {
		t.Error("error of CopyWithOption time to int64 by unixmill")
	}
	copier.CopyWithOption(&c, &now, copier.Option{TimeFormat: "unix"})
	if now.Unix() != c {
		t.Error("error of CopyWithOption time to int64 by unixmill")
	}
	var timstreap int64 = 165000760462
	b = time.Time{}
	copier.CopyWithOption(&b, &timstreap, copier.Option{TimeFormat: "unixmill"})
	// fmt.Println(b)
}

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

func checkTimeData(data1 TimInt, data3 TimTim, t *testing.T) {
	if data1.Time1 != data3.Time1.Unix() {
		t.Error("data1.Time1 != data3.Time1.Unix()")
	}
	if data1.Time2 != 0 {
		t.Error("data1.Time2 != 0")
	}
	if data1.Time3 != data3.Time3.Unix() {
		t.Error("data1.Time3 != data3.Time3.Unix()")
	}
	if data1.Time4 != data3.Time4.UnixMilli() {
		fmt.Println(data1.Time4, data3.Time4.UnixMilli())
		t.Error("data1.Time4 != data3.Time4.UnixMilli()")
	}
	if data1.Time5 != nil {
		t.Errorf("data1.Time5 != nil")
	}
	if data1.Time6 != nil {
		t.Errorf("data1.Time6 != nil")
	}
	if data1.Time7 != 0 && !data3.Time7.IsZero() {
		t.Error("data1.Time7 != 0 && !data3.Time7.IsZero()")
	}
}

func TestCopyTimeAndInt642(t *testing.T) {
	now := time.Now()
	var tim = TimTim{
		Time1: now,
		Time2: now,
		Time3: &now,
		Time4: &now,
	}
	var num TimInt
	if err := copier.Copy(&num, &tim); err != nil {
		t.Errorf("Copy run error in struct from time to int64: %s", err)
	}
	checkTimeData(num, tim, t)

	tim = TimTim{}
	if err := copier.Copy(&tim, &num); err != nil {
		t.Errorf("Copy run error in struct from int64 to time: %s", err)
	}
	checkTimeData(num, tim, t)
}
