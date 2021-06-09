package main

import (
	"encoding/json"
	"fmt"
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

func GetConfiguration() Configuration {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configuration Configuration
	_ = json.Unmarshal(byteValue, &configuration)
	return configuration
}

type Configuration struct {
	DbUser string `json:"db-user"`
	DbPassword string `json:"db-password"`
	DbPort string `json:"db-port"`
}

type ConfigurationFile struct {
	Name          string        `json:"name"`
	Author        string        `json:"author"`
	Configuration Configuration `json:"configuration"`
}