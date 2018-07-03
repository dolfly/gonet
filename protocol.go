package gonet

import (
	net "github.com/dolfly/gonet/net"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn net.Conn) (Packet, error)
}
