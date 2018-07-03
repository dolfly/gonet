package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/dolfly/gonet"
	"github.com/dolfly/gonet/cmd/echo"
	net "github.com/dolfly/gonet/net"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gonet.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gonet.Conn, p gonet.Packet) bool {
	echoPacket := p.(*echo.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(echo.NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (this *Callback) OnClose(c *gonet.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	listener, err := net.ListenUDP("0.0.0.0", 7800)
	checkError(err)

	// creates a server
	config := &gonet.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gonet.NewServer(listener, config, &Callback{}, &echo.EchoProtocol{})
	// starts service
	go srv.Start()
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
