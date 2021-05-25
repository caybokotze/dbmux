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
var Bsize uint
var Verbose bool
var Dbh *sql.DB

func main() {
	var bind, backend, logTo string
	var buffer uint
	var verbose bool
	var conf string

	flag.StringVar(&bind, "bind", ":8002", "locate ip and port")
	flag.StringVar(&backend, "backend", "127.0.0.1:8003", "backend server ip and port")
	flag.StringVar(&logTo, "logTo", "stdout", "stdout or syslog")
	flag.UintVar(&buffer, "buffer", 4096, "buffer size")
	flag.BoolVar(&verbose, "verbose", false, "print verbose sql query")
	flag.StringVar(&conf, "conf", "", "config file to verify database and record sql query")
	flag.Parse()
	Bsize = buffer
	Verbose = verbose

	conf_fh, err := get_config(conf)
	if err != nil {
		log.Printf("Can't get config info, skip insert log to mysql...\n")
	} else {
		backendDsn, _ := get_backend_dsn(conf_fh)
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
