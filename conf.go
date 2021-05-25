package main

import (
	"github.com/arstercz/goconfig"
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