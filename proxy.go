package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Proxy struct {
	proxyTcp, databaseTcp *net.TCPAddr
	sessionsCount     int32
	pool              *recycler
}

func CreateNewProxy(proxyPort, databasePort uint, size uint32) *Proxy {
	proxyTcpAddress, err := net.ResolveTCPAddr("tcp", strconv.Itoa(int(proxyPort)))
	if err != nil {
		log.Fatalln("resolve proxyTcp error:", err)
	}

	databaseTcpAddress, err := net.ResolveTCPAddr("tcp", strconv.Itoa(int(databasePort)))
	if err != nil {
		log.Fatalln("resolve backend error:", err)
	}

	return &Proxy{
		proxyTcp: proxyTcpAddress,
		databaseTcp: databaseTcpAddress,
		sessionsCount: 0,
		pool: NewRecycler(size),
	}
}

/* - Proxy struct helpers - */

func (t *Proxy) pipeTCPConnection(dst, src *Conn, c chan int64, tag string) {
	defer func() {
		dst.CloseWrite()
		dst.CloseRead()
	}()
	if strings.EqualFold(tag, "send") {
		ProxyLog(src, dst)
		c <- 0
	} else {
		n, err := io.Copy(dst, src)
		if err != nil {
			log.Print(err)
		}
		c <- n
	}
}

func (t *Proxy) transport(conn net.Conn) {
	start := time.Now()
	conn2, err := net.DialTCP("tcp", nil, t.databaseTcp)
	if err != nil {
		log.Print(err)
		return
	}
	connectTime := time.Now().Sub(start)
	log.Printf("proxy: %s ==> %s", conn2.LocalAddr().String(),
		conn2.RemoteAddr().String())
	start = time.Now()
	readChan := make(chan int64)
	writeChan := make(chan int64)
	var readBytes, writeBytes int64

	atomic.AddInt32(&t.sessionsCount, 1)
	var bindConn, backendConn *Conn
	bindConn = NewConn(conn, t.pool)
	backendConn = NewConn(conn2, t.pool)

	go t.pipeTCPConnection(backendConn, bindConn, writeChan, "send")
	go t.pipeTCPConnection(bindConn, backendConn, readChan, "receive")

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
