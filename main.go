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

func Initialise() {
	hostPort := flag.Uint("host", 3307, "Specify the port the sql server is running on")
	proxyPort := flag.Uint("proxy", 3308, "Specify the port where the current server instance is running")
	verbosityEnabled := flag.Bool( "enable-verbosity", true, "Enable verbosity to see the output in terminal")
	flag.Parse()
	log.SetOutput(os.Stdout)
	conf, err := fetchConfiguration(hostPort, proxyPort)
	if err != nil {
		log.Fatal("Configuration could not be found for this service, please make sure you have a valid config file.")
	}
	dbn, dbErr := database.CreateConnectionToDbHost(conf)
	if dbErr != nil {
		log.Fatal("Count not create a connection to the database")
	}

	p := proxy.CreateNewProxy(proxy.Arguments{
		ProxyPort:      *proxyPort,
		HostPort:       *hostPort,
		BufferSize:     4096,
		ThreadPoolSize: 50,
		VerbosityEnabled: *verbosityEnabled,
		DatabaseHost: dbn,
	})

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

/*
Waits for the TCP channel to open or close.
 */
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
