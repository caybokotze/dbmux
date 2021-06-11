package main

import (
	"net"
	"time"
)

type Connection struct {
	connection net.Conn
	pool       *recycler
}

func NewTcpConnection(conn net.Conn, pool *recycler) *Connection {
	return &Connection{
		connection: conn,
		pool:       pool,
	}
}

func (c *Connection) Read(b []byte) (int, error) {
	_ = c.connection.SetReadDeadline(time.Now().Add(30 * time.Minute))
	n, err := c.connection.Read(b)
	return n, err
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.connection.Write(b)
}

func (c *Connection) Close() {
	_ = c.connection.Close()
}

func (c *Connection) CloseRead() {
	if conn, ok := c.connection.(*net.TCPConn); ok {
		_ = conn.CloseRead()
	}
}

func (c *Connection) CloseWrite() {
	if conn, ok := c.connection.(*net.TCPConn); ok {
		_ = conn.CloseWrite()
	}
}
