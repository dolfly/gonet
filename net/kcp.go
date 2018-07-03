package net

import (
	"fmt"
	"net"

	"github.com/dolfly/gonet/log"

	kcp "github.com/xtaci/kcp-go"
)

type KcpListener struct {
	addr      net.Addr
	listener  net.Listener
	accept    chan Conn
	closeFlag bool
	log.Logger
}

func ListenKCP(bindAddr string, bindPort int) (l *KcpListener, err error) {
	listener, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", bindAddr, bindPort), nil, 10, 3)
	if err != nil {
		return l, err
	}
	listener.SetReadBuffer(4194304)
	listener.SetWriteBuffer(4194304)

	l = &KcpListener{
		addr:      listener.Addr(),
		listener:  listener,
		accept:    make(chan Conn),
		closeFlag: false,
		Logger:    log.NewPrefixLogger(""),
	}

	go func() {
		for {
			conn, err := listener.AcceptKCP()
			if err != nil {
				if l.closeFlag {
					close(l.accept)
					return
				}
				continue
			}
			conn.SetStreamMode(true)
			conn.SetWriteDelay(true)
			conn.SetNoDelay(1, 20, 2, 1)
			conn.SetMtu(1350)
			conn.SetWindowSize(1024, 1024)
			conn.SetACKNoDelay(false)

			l.accept <- WrapConn(conn)
		}
	}()
	return l, err
}

func (l *KcpListener) Accept() (Conn, error) {
	conn, ok := <-l.accept
	if !ok {
		return conn, fmt.Errorf("channel for kcp listener closed")
	}
	return conn, nil
}

func (l *KcpListener) Close() error {
	if !l.closeFlag {
		l.closeFlag = true
		l.listener.Close()
	}
	return nil
}

func (l *KcpListener) Addr() net.Addr {
	return l.addr
}

func ConnectKcpServer(addr string) (c Conn, err error) {
	kcpConn, errRet := kcp.DialWithOptions(addr, nil, 10, 3)
	if errRet != nil {
		err = errRet
		return
	}
	kcpConn.SetStreamMode(true)
	kcpConn.SetWriteDelay(true)
	kcpConn.SetNoDelay(1, 20, 2, 1)
	kcpConn.SetWindowSize(128, 512)
	kcpConn.SetMtu(1350)
	kcpConn.SetACKNoDelay(false)
	kcpConn.SetReadBuffer(4194304)
	kcpConn.SetWriteBuffer(4194304)
	c = WrapConn(kcpConn)
	return
}

/*
func NewKcpConnFromUdp(conn *net.UDPConn, connected bool, raddr string) (net.Conn, error) {
	kcpConn, err := kcp.NewConnEx(1, connected, raddr, nil, 10, 3, conn)
	if err != nil {
		return nil, err
	}
	kcpConn.SetStreamMode(true)
	kcpConn.SetWriteDelay(true)
	kcpConn.SetNoDelay(1, 20, 2, 1)
	kcpConn.SetMtu(1350)
	kcpConn.SetWindowSize(1024, 1024)
	kcpConn.SetACKNoDelay(false)
	return kcpConn, nil
}*/
