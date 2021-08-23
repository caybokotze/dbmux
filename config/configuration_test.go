package config

import "testing"

func TestGetConfiguration(t *testing.T) {
	conf, err := GetConfiguration()
	if err != nil {
		t.Fatal("This test failed: ", err)
	}
	if conf.VerbosityEnabled == true || conf.VerbosityEnabled == false {
		t.Log("Verbosity is valid")
	}
}

type comparator interface {
	or(compare, to bool) bool
	equal(compare, to bool) bool
}

type values struct {
	compare bool
	to bool
}

func (v *values) equal(compare, to bool) bool{
	if v.compare == v.to {
		return true
	}
	return false
}

func (v *values) or(compare, to bool) bool{
	if v.compare == v.to {
		return true
	}
	return false
}