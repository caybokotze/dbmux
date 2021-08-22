package proxy

import (
	"github.com/caybokotze/dbmux/logging"
	"github.com/caybokotze/dbmux/tcp"
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
	sessionsCount         int32
	pool                  *Recycler
	bufferSize            uint
}

type Arguments struct {
	ProxyPort      uint
	HostPort       uint
	BufferSize     uint
	ThreadPoolSize uint32
	VerbosityEnabled bool
}

func CreateNewProxy(arguments Arguments) *Proxy {
	proxyTcpAddress, err := net.ResolveTCPAddr("tcp", strconv.Itoa(int(arguments.ProxyPort)))
	if err != nil {
		log.Fatalln("resolve proxyTcp error:", err)
	}

	databaseTcpAddress, err := net.ResolveTCPAddr("tcp", strconv.Itoa(int(arguments.HostPort)))
	if err != nil {
		log.Fatalln("resolve backend error:", err)
	}

	return &Proxy{
		proxyTcp:      proxyTcpAddress,
		databaseTcp:   databaseTcpAddress,
		sessionsCount: 0,
		pool:          newRecycler(arguments.ThreadPoolSize),
		bufferSize:    arguments.BufferSize,
	}
}

/* - Proxy struct helpers - */

func (t *Proxy) pipeTCPConnection(dst, src *tcp.Connection, c chan int64, tag string) {
	defer func() {
		dst.CloseWrite()
		dst.CloseRead()
	}()
	if strings.EqualFold(tag, "send") {
		logging.ProxyLog(src, dst, t.bufferSize)
		c <- 0
	} else {
		n, err := io.Copy(dst, src)
		if err != nil {
			log.Print(err)
		}
		c <- n
	}
}

func (t *Proxy) transport(proxy net.Conn) {
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
	var serverConnection, proxyConnection *tcp.Connection
	serverConnection = tcp.NewTcpConnection(proxy, t.pool)
	proxyConnection = tcp.NewTcpConnection(databaseConnection, t.pool)

	go t.pipeTCPConnection(serverConnection, proxyConnection, writeChan, "send")
	go t.pipeTCPConnection(proxyConnection, serverConnection, readChan, "receive")

	readBytes = <-readChan
	writeBytes = <-writeChan
	transferTime := time.Now().Sub(start)
	log.Printf("r: %d w:%d ct:%.3f t:%.3f [#%d]", readBytes, writeBytes,
		connectTime.Seconds(), transferTime.Seconds(), t.sessionsCount)
	atomic.AddInt32(&t.sessionsCount, -1)
}

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
		go t.transport(conn)
	}
}
