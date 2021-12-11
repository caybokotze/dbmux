package proxy

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Proxy struct {
	proxyTcp,
	databaseTcp *net.TCPAddr
	databaseHost *sql.DB
	sessionsCount         int32
	pool                  *Recycler
	bufferSize            uint
	verbosity bool
}

type Arguments struct {
	ProxyPort      uint
	HostPort       uint
	BufferSize     uint
	ThreadPoolSize uint32
	VerbosityEnabled bool
	DatabaseHost *sql.DB
}

/*
-------------------------------------------------------------------------------------------------------
This function resolves for the running mysql connection as well as the proxy port specified.
This is required to determine whether something is making use of the port or not.
-------------------------------------------------------------------------------------------------------
 */
func CreateNewProxy(arguments Arguments) *Proxy {
	var proxyPort = fmt.Sprintf("%s%s", ":", strconv.Itoa(int(arguments.ProxyPort)))
	proxyTcpAddress, err := net.ResolveTCPAddr("tcp", proxyPort)
	if err != nil {
		log.Fatalln("Could not resolve the host port:", err)
	}
	log.Printf("Proxy running on port: %s", strconv.Itoa(int(arguments.ProxyPort)))

	var hostPort = fmt.Sprintf("%s%s", ":", strconv.Itoa(int(arguments.HostPort)))
	databaseTcpAddress, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		log.Fatalln("could not resolve the proxy port:", err)
	}
	log.Printf("Connected to host on port: %s", strconv.Itoa(int(arguments.HostPort)))

	return &Proxy{
		proxyTcp:      proxyTcpAddress,
		databaseTcp:   databaseTcpAddress,
		sessionsCount: 0,
		pool:          newRecycler(arguments.ThreadPoolSize),
		bufferSize:    arguments.BufferSize,
		verbosity: arguments.VerbosityEnabled,
		databaseHost: arguments.DatabaseHost,
	}
}

/* - Proxy struct helpers - */

func (t *Proxy) pipeTCPConnection(dst, src *Connection, c chan int64, tag string) {
	defer func() {
		dst.CloseWrite()
		dst.CloseRead()
	}()
	if strings.EqualFold(tag, "send") {
		Log(LogConfiguration{
			Source: src,
			Destination: dst,
			BufferSize: t.bufferSize,
			Verbosity: t.verbosity,
			DatabaseHost: t.databaseHost,
		})
		c <- 0
	} else {
		n, err := io.Copy(dst, src)
		if err != nil {
			log.Print(err)
		}
		c <- n
	}
}

func (t *Proxy) transportTCPConnection(proxy net.Conn) {
	start := time.Now()
	databaseConnection, err := net.DialTCP("tcp", nil, t.databaseTcp)
	if err != nil {
		log.Print(err)
		return
	}
	connectTime := time.Now().Sub(start)
	log.Printf("proxy: %s ==> %s", databaseConnection.LocalAddr().String(),
		databaseConnection.RemoteAddr().String())
	start = time.Now()
	readChan := make(chan int64)
	writeChan := make(chan int64)
	var readBytes, writeBytes int64

	atomic.AddInt32(&t.sessionsCount, 1)
	var serverConnection, proxyConnection *Connection
	serverConnection = NewTcpConnection(proxy, t.pool)
	proxyConnection = NewTcpConnection(databaseConnection, t.pool)

	go t.pipeTCPConnection(proxyConnection, serverConnection, writeChan, "send")
	go t.pipeTCPConnection(serverConnection, proxyConnection, readChan, "receive")

	readBytes = <-readChan
	writeBytes = <-writeChan
	transferTime := time.Now().Sub(start)
	log.Printf("r: %d w:%d ct:%.3f t:%.3f [#%d]", readBytes, writeBytes,
		connectTime.Seconds(), transferTime.Seconds(), t.sessionsCount)
	atomic.AddInt32(&t.sessionsCount, -1)
}

/*
Checks that the TCP port is not being used by another application before transporting the TCP connection.
 */
func (t *Proxy) StartTcpProxying() {
	ln, err := net.ListenTCP("tcp", t.proxyTcp)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		log.Printf("client: %s ==> %s", conn.RemoteAddr().String(),
			conn.LocalAddr().String())
		go t.transportTCPConnection(conn)
	}
}
