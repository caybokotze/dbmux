package main

import (
	"net"
	"time"
)

type Conn struct {
	conn net.Conn
	pool *recycler
}

func NewTcpConnection(conn net.Conn, pool *recycler) *Conn {
	return &Conn{
		conn: conn,
		pool: pool,
	}
}

func (c *Conn) Read(b []byte) (int, error) {
	_ = c.conn.SetReadDeadline(time.Now().Add(30 * time.Minute))
	n, err := c.conn.Read(b)
	return n, err
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Conn) Close() {
	_ = c.conn.Close()
}

func (c *Conn) CloseRead() {
	if conn, ok := c.conn.(*net.TCPConn); ok {
		_ = conn.CloseRead()
	}
}

func (c *Conn) CloseWrite() {
	if conn, ok := c.conn.(*net.TCPConn); ok {
		_ = conn.CloseWrite()
	}
}
