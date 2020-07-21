package frontier

import "net"

type Driver interface {
	OnInit() error
	OnOpen(nConn net.Conn) (Conn,error)
	OnMessage(fConn Conn,message []byte)
	OnClose(fConn Conn)
}
