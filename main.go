package main

import (
	_ "bytes"
	_ "database/sql"
	"flag"
	"github.com/caybokotze/dbmux/config"
	"github.com/caybokotze/dbmux/database"
	"github.com/caybokotze/dbmux/proxy"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
	conf, err := fetchConfiguration(bindingPort, proxyPort)
	if err != nil {
		log.Fatal("Configuration could not be found for this service, please make sure you have a valid config file.")
	}

	_, err = database.CreateConnectionToDbHost(conf)

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

func fetchConfiguration(bindingPort, proxyPort *uint)(configuration config.Configuration, err error) {
	conf, err := config.GetConfiguration()
	if err != nil {
		log.Println("Error fetching config")
		return config.Configuration{}, err
	}
	if *proxyPort != 0 {
		conf.ProxyPort = *proxyPort
	}
	if *bindingPort == 0 {
		conf.DbPort = *bindingPort
	}
	return conf, nil
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
