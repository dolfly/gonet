package net

import (
	"fmt"
	"net"
	"time"

	"github.com/dolfly/gonet/log"
)

type TcpListener struct {
	addr      net.Addr
	listener  net.Listener
	accept    chan Conn
	closeFlag bool
	log.Logger
}

func ListenTCP(bindAddr string, bindPort int) (l *TcpListener, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		return l, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return l, err
	}

	//listener.SetDeadline(time.Now().Add(20 * time.Second))

	l = &TcpListener{
		addr:      listener.Addr(),
		listener:  listener,
		accept:    make(chan Conn),
		closeFlag: false,
		Logger:    log.NewPrefixLogger(""),
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				if l.closeFlag {
					close(l.accept)
					return
				}
				continue
			}

			c := NewTcpConn(conn)
			l.accept <- c
		}
	}()
	return l, err
}

// Wait util get one new connection or listener is closed
// if listener is closed, err returned.
func (l *TcpListener) Accept() (Conn, error) {
	conn, ok := <-l.accept
	if !ok {
		return conn, fmt.Errorf("channel for tcp listener closed")
	}
	return conn, nil
}

func (l *TcpListener) Close() error {
	if !l.closeFlag {
		l.closeFlag = true
		l.listener.Close()
	}
	return nil
}

func (l *TcpListener) Addr() net.Addr {
	return l.addr
}

// Wrap for TCPConn.
type TcpConn struct {
	net.Conn
	log.Logger
}

func NewTcpConn(conn net.Conn) (c *TcpConn) {
	c = &TcpConn{
		Conn:   conn,
		Logger: log.NewPrefixLogger("TCP"),
	}
	return
}

func ConnectTcpServer(addr string) (c Conn, err error) {
	servertAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, servertAddr)
	if err != nil {
		return
	}
	c = NewTcpConn(conn)
	return
}
