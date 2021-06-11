package main

import (
	"encoding/json"
	"github.com/arstercz/goconfig"
	"io/ioutil"
	"os"
)

func getConfig(conf string) (c *goconfig.ConfigFile, err error) {
	c, err = goconfig.ReadConfigFile(conf)
	if err != nil {
		return c, err
	}
	return c, nil
}

func getBackendDsn(c *goconfig.ConfigFile) (dsn string, err error) {
	dsn, err = c.GetString("backend", "dsn")
	if err != nil {
		return dsn, err
	}
	return dsn, nil
}

func GetConfiguration() (config Configuration, err error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return Configuration{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configuration Configuration
	_ = json.Unmarshal(byteValue, &configuration)
	return configuration, nil
}

type Configuration struct {
	DbUser string `json:"db-user"`
	DbPassword string `json:"db-password"`
	DbPort uint `json:"db-port"`
	ProxyPort uint `json:"proxy-port"`
	DbBuffer uint `json:"db-buffer"`
}

type ConfigurationFile struct {
	Name          string        `json:"name"`
	Author        string        `json:"author"`
	Configuration Configuration `json:"configuration"`
}