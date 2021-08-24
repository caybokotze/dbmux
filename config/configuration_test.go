package config

import (
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	conf, err := GetConfiguration()
	if err != nil {
		//var result = tkit.Compare(false).To(false)
		t.Fatal("This test failed: ", err)
	}
	if conf.VerbosityEnabled == true || conf.VerbosityEnabled == false {
		t.Log("Verbosity is valid")
	}
}