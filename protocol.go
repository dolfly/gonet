package gonet

import (
	"net"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn *conn.Conn) (Packet, error)
}
