package proxy

import (
	"net"
	"time"
)

type Connection struct {
	Connection net.Conn
	pool       *Recycler
}

func NewTcpConnection(conn net.Conn, pool *Recycler) *Connection {
	return &Connection{
		Connection: conn,
		pool:       pool,
	}
}

func (c *Connection) Read(b []byte) (int, error) {
	_ = c.Connection.SetReadDeadline(time.Now().Add(30 * time.Minute))
	n, err := c.Connection.Read(b)
	return n, err
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.Connection.Write(b)
}

func (c *Connection) Close() {
	_ = c.Connection.Close()
}

func (c *Connection) CloseRead() {
	if conn, ok := c.Connection.(*net.TCPConn); ok {
		_ = conn.CloseRead()
	}
}

func (c *Connection) CloseWrite() {
	if conn, ok := c.Connection.(*net.TCPConn); ok {
		_ = conn.CloseWrite()
	}
}
