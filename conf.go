package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetConfiguration() (config Configuration, err error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return Configuration{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configurationFile ConfigurationFile
	_ = json.Unmarshal(byteValue, &configurationFile)
	return configurationFile.Configuration, nil
}

type Configuration struct {
	DbUser string `json:"db-user"`
	DbPassword string `json:"db-password"`
	DbPort uint `json:"db-port"`
	DbHost string `json:"db-host"`
	ProxyPort uint `json:"proxy-port"`
	DbBuffer uint `json:"db-buffer"`
	DbSchema string `json:"db-schema"`
}

type ConfigurationFile struct {
	Name          string        `json:"name"`
	Author        string        `json:"author"`
	Configuration Configuration `json:"configuration"`
}