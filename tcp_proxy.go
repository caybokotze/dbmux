package main

import (
	_ "bytes"
	"database/sql"
	_ "database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeout = time.Second * 2

func main() {
	Initialise()
}

var DatabaseHost *sql.DB
var VerbosityEnabled = false

func Initialise() {
	bindingPort := flag.Uint("bind-to", 3601, "Specify the port you will be accessing from")
	proxyPort := flag.Uint("proxy-to", 3600, "Specify the port where the current server instance is running")
	flag.BoolVar(&VerbosityEnabled, "enable-verbosity", false, "Select whether or not verbosity is enabled")
	flag.Parse()

	log.SetOutput(os.Stdout)
	config, err := fetchConfiguration(bindingPort, proxyPort)
	if err != nil {
		log.Fatal("Configuration could not be found for this service, please make sure you have a valid configuration file.")
	}

	DatabaseHost, err = databaseHost(config)
	if err != nil {
		log.Fatal("Count not create a connection to the database")
	}

	p := CreateNewProxy(
		config.ProxyPort,
		config.DbPort,
		uint32(config.DbBuffer))
	log.Println("portproxy started.")
	go p.StartTcpProxying()
	waitForSignal()
}

func fetchConfiguration(bindingPort, proxyPort *uint) (configuration Configuration, err error) {
	config, err := GetConfiguration()
	if err != nil {
		log.Println("Error fetching configuration")
		return Configuration{}, err
	}
	if *proxyPort != 0 {
		config.ProxyPort = *proxyPort
	}
	if *bindingPort == 0 {
		config.DbPort = *bindingPort
	}
	return config, nil
}

func waitForSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan)
	for sig := range sigChan {
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			log.Printf("terminated by signal %v\n", sig)
		} else {
			log.Printf("received signal: %v, ignore\n", sig)
		}
	}
}
