package main

import (
	"fmt"
	"log"

	"time"

	"github.com/dolfly/gonet/cmd/echo"
	"github.com/dolfly/gonet/net"
)

func main() {
	conn, err := net.ConnectServer("udp", "127.0.0.1:7800")
	checkError(err)

	echoProtocol := &echo.EchoProtocol{}

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		conn.Write(echo.NewEchoPacket([]byte("hello"), false).Serialize())

		// read
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*echo.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
