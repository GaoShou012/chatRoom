package frontier

import "net"

type Conn interface {
	GetConn() net.Conn
	GetId() string
	Writer(message []byte) error
	Reader() ([]byte,error)
}
