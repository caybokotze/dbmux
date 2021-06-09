package main

import (
	"bytes"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeout = time.Second * 2
var BSize uint
var Verbose bool
var Dbh *sql.DB

func main() {
	Initialise()
}

func Initialise() {
	bindingPort := flag.Uint("bind-to", 3602, "Specify the port you will be accessing from")
	proxyPort := flag.Uint("proxy-to", 3600, "Specify the port where the current server instance is running")
	verbosity := flag.Bool("enable-verbosity", false, "Select whether or not verbosity is enabled")

	flag.Parse()
	BSize = buffer
	Verbose = verbose

	confFh, err := getConfig(conf)
	if err != nil {
		log.Printf("Can't get config info, skip insert log to mysql...\n")
	} else {
		backendDsn, _ := getBackendDsn(confFh)
		Dbh, err = dbh(backendDsn)
		if err != nil {
			log.Printf("Can't get database handle, skip insert log to mysql...\n")
		}
		defer Dbh.Close()
	}

	log.SetOutput(os.Stdout)
	if logTo == "syslog" {
		var (
			buf    bytes.Buffer
			logger = log.New(&buf, "INFO: ", log.Lshortfile)
			infoF  = func(info string) {
				_ = logger.Output(2, info)
			}
		)
		infoF("port proxying...")
	}

	p := New(bind, backend, uint32(buffer))
	log.Println("portproxy started.")
	go p.Start()
	waitSignal()
}

func waitSignal() {
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
