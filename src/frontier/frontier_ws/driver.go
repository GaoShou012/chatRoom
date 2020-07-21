package frontier_ws

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/src/frontier"
	"net"
	"net/url"
)

var _ frontier.Driver = &Driver{}

type Driver struct {
	Room
}

func (d *Driver) OnInit() error {
	fmt.Println("on init")
	d.Room.Init()
	return nil
}

func (d *Driver) OnOpen(nConn net.Conn) (frontier.Conn, error) {
	token := ""
	conn := &Conn{conn: nConn}
	_, err := ws.Upgrader{
		OnRequest: func(uri []byte) error {
			u, err := url.Parse(string(uri))
			if err != nil {
				return err
			}
			if t := u.Query().Get("token"); t != "" {
				token = t
			}
			return nil
		},
		OnHeader: func(key, value []byte) error {
			return nil
		},
		OnBeforeUpgrade: func() (_ ws.HandshakeHeader, err error) {
			// 返回错误，连接升级协议失败
			// 返回nil,连接升级协议成功
			d.Room.Join(conn)
			return
		},
	}.Upgrade(conn.conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *Driver) OnMessage(fConn frontier.Conn, message []byte) {
	conn := fConn.(*Conn)
	msg := Message{
		RuledOut: make(map[*Conn]bool),
		Body:     message,
	}
	msg.RuledOut[conn] = true

	d.Room.Broadcast(msg)
}

func (d *Driver) OnClose(fConn frontier.Conn) {
	fmt.Println("i am onClose", fConn)
}
