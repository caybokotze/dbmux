package configuration

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
	var configurationFile File
	_ = json.Unmarshal(byteValue, &configurationFile)
	return configurationFile.Configuration, nil
}

type Configuration struct {
	DbUser           string `json:"db-user"`
	DbPassword       string `json:"db-password"`
	DbPort           uint   `json:"db-port"`
	DbHostIp         string `json:"db-host"`
	ProxyPort        uint   `json:"proxy-port"`
	DbBuffer         uint   `json:"db-buffer"`
	DbSchema         string `json:"db-schema"`
	VerbosityEnabled bool   `json:"verbosity"`
}

type File struct {
	Name          string        `json:"name"`
	Author        string        `json:"author"`
	Configuration Configuration `json:"proxy-configuration"`
}