package main

import (
	_ "bytes"
	_ "database/sql"
	"flag"
	"github.com/caybokotze/dbmux/configuration"
	"github.com/caybokotze/dbmux/database"
	"github.com/caybokotze/dbmux/proxy"
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

var verbosityEnabled = false

func Initialise() {
	bindingPort := flag.Uint("bind-to", 3306, "Specify the port the sql server is running on")
	proxyPort := flag.Uint("proxy-to", 3600, "Specify the port where the current server instance is running")
	flag.BoolVar(&verbosityEnabled, "enable-verbosity", false, "Enable verbosity to see the output in terminal")
	flag.Parse()

	log.SetOutput(os.Stdout)
	config, err := fetchConfiguration(bindingPort, proxyPort)
	if err != nil {
		log.Fatal("Configuration could not be found for this service, please make sure you have a valid configuration file.")
	}

	_, err = database.CreateConnectionToDbHost(config)

	if err != nil {
		log.Fatal("Count not create a connection to the database")
	}

	p := proxy.CreateNewProxy(proxy.Arguments{
		ProxyPort:      *proxyPort,
		HostPort:       *bindingPort,
		BufferSize:     0,
		ThreadPoolSize: 50,
		VerbosityEnabled: verbosityEnabled,
	})

	log.Println("portproxy started.")
	go p.StartTcpProxying()
	waitForSignal()
}

func fetchConfiguration(bindingPort, proxyPort *uint)(conf configuration.Configuration, err error) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println("Error fetching configuration")
		return configuration.Configuration{}, err
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
