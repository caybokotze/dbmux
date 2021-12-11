package config

import (
	tk "github.com/caybokotze/go-testing-kit"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	conf, err := GetConfiguration()
	if err != nil {
		t.Fatal("This test failed: ", err)
	}
	if tk.Compare(conf.VerbosityEnabled).To(false) {
		t.Log("Verbosity is valid")
	}
}

func TestGetConfiguration_GetDbBuffer_ShouldNotBeNull(t *testing.T) {
	conf, err := GetConfiguration()
	if err != nil {
		t.Fail()
	}
	if tk.Compare(conf.DbBuffer).To(0) {
		t.Logf("value doesn't exist...: %d", conf.DbBuffer)
		t.Fail()
	}
	if tk.Compare(conf.DbSchema).To("") {
		t.Log("Db Schema is not set")
		t.Fail()
	}
	if tk.Compare(conf.DbUser).To("") {
		t.Log("Db user should be set")
		t.Fail()
	}
	if tk.Compare(conf.DbPassword).To("") {
		t.Log("Db password should be set")
		t.Fail()
	}
	if tk.Compare(conf.DbBuffer).To(0) {
		t.Log("Db buffer should be set")
		t.Fail()
	}
	if tk.Compare(conf.DbPort).To(0) {
		t.Log("Db port should be set")
		t.Fail()
	}
	if tk.Compare(conf.ProxyPort).To(0) {
		t.Log("Proxy port should be set")
		t.Fail()
	}
	if tk.Compare(conf.ThreadPoolCount).To(0) {
		t.Log("Thread pool count should be set")
		t.Fail()
	}
}