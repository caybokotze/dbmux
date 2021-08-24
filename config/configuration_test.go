package config

import (
	kit "github.com/caybokotze/go-testing-kit"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	conf, err := GetConfiguration()
	if err != nil {
		t.Fatal("This test failed: ", err)
	}
	if conf.VerbosityEnabled == true || conf.VerbosityEnabled == false {
		t.Log("Verbosity is valid")
	}
}