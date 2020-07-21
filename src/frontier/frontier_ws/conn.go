package frontier_ws

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/src/frontier"
	"io/ioutil"
	"net"
	"time"
)

var _ frontier.Conn = &Conn{}
var (
	ReadMessageTimeout  = 100 * time.Millisecond
	WriteMessageTimeout = 100 * time.Millisecond
)

type Conn struct {
	conn net.Conn
}

func (c Conn) GetConn() net.Conn {
	return c.conn
}

func (c Conn) GetId() string {
	panic("implement me")
}

func (c Conn) Writer(message []byte) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(WriteMessageTimeout)); err != nil {
		return err
	}
	w := wsutil.NewWriter(c.conn, ws.StateServerSide, ws.OpText)
	_,err := w.Write(message)
	if err != nil {
		return err
	}
	return w.Flush()
}

func (c Conn) Reader() ([]byte, error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(ReadMessageTimeout)); err != nil {
		return nil, err
	}
	h, r, err := wsutil.NextReader(c.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(c.conn, ws.StateServerSide)(h, r)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}



