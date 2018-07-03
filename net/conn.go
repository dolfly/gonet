package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/dolfly/gonet/log"
)

type Conn interface {
	net.Conn
}
type WrapLogConn struct {
	net.Conn
	log.Logger
}

func WrapConn(c net.Conn) Conn {
	return &WrapLogConn{
		Conn:   c,
		Logger: log.NewPrefixLogger("WRAP"),
	}
}

type WrapReadWriteCloserConn struct {
	io.ReadWriteCloser
	log.Logger

	underConn net.Conn
}

func WrapReadWriteCloserToConn(rwc io.ReadWriteCloser, underConn net.Conn) Conn {
	return &WrapReadWriteCloserConn{
		ReadWriteCloser: rwc,
		Logger:          log.NewPrefixLogger(""),
		underConn:       underConn,
	}
}

func (conn *WrapReadWriteCloserConn) LocalAddr() net.Addr {
	if conn.underConn != nil {
		return conn.underConn.LocalAddr()
	}
	return (*net.TCPAddr)(nil)
}

func (conn *WrapReadWriteCloserConn) RemoteAddr() net.Addr {
	if conn.underConn != nil {
		return conn.underConn.RemoteAddr()
	}
	return (*net.TCPAddr)(nil)
}

func (conn *WrapReadWriteCloserConn) SetDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetDeadline(t)
	}
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (conn *WrapReadWriteCloserConn) SetReadDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetReadDeadline(t)
	}
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (conn *WrapReadWriteCloserConn) SetWriteDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetWriteDeadline(t)
	}
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func ConnectServer(protocol string, addr string) (c Conn, err error) {
	switch protocol {
	case "tcp":
		return ConnectTcpServer(addr)
	case "kcp":
		return ConnectKcpServer(addr)
	case "udp":
		return ConnectUdpServer(addr)
	default:
		return nil, fmt.Errorf("unsupport protocol: %s", protocol)
	}
}
