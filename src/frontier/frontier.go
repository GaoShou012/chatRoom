package frontier

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/mailru/easygo/netpoll"
	"net"
	"time"
)

func NewFrontier(ln net.Listener,driver Driver) (*Frontier, error) {
	poller, err := netpoll.New(nil)
	if err != nil {
		return nil, err
	}
	desc := netpoll.Must(netpoll.HandleListener(ln, netpoll.EventRead|netpoll.EventOneShot))

	f := &Frontier{
		ln:     ln,
		accept: make(chan error, 1),
		desc:   desc,
		poller: poller,
		Driver: driver,
	}

	if err := driver.OnInit(); err != nil {
		return nil,err
	}

	return f, nil
}

type Frontier struct {
	ln     net.Listener
	accept chan error

	desc   *netpoll.Desc
	poller netpoll.Poller

	Driver
}

func (f *Frontier) Start() error {
	return f.poller.Start(f.desc, func(event netpoll.Event) {
		go func() {
			conn, err := f.ln.Accept()
			if err != nil {
				f.accept <- err
				return
			}
			f.accept <- nil
			f.onOpen(conn)
		}()
		err := <-f.accept
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Errorf("accept error:%v\n", err)
				time.Sleep(time.Millisecond * 5)
			}
		}
		f.poller.Resume(f.desc)
	})
}

func (f *Frontier) Stop() error {
	if err := f.poller.Stop(f.desc); err != nil {
		return err
	}
	return f.ln.Close()
}

func (f *Frontier) onOpen(nConn net.Conn) {
	// 从驱动层获取连接
	driverConn, err := f.Driver.OnOpen(nConn)
	if err != nil {
		glog.Errorln(err)
		return
	}

	// 处理连接事件
	// 接受数据、连接断开事件、ping等等
	desc := netpoll.Must(netpoll.HandleRead(nConn))
	f.poller.Start(desc, func(event netpoll.Event) {
		if event&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			glog.Infoln("%s: receive: %v; close connection", nConn, event)

			f.poller.Stop(desc)
			f.Driver.OnClose(driverConn)
			return
		}

		// Here we can read some new message from connection.
		// We can not read it right here in callback, because then we will
		// block the poller's inner loop.
		// We do not want to spawn a new goroutine to read single message.
		// But we want to reuse previously spawned goroutine.
		go func() {
			msg, err := driverConn.Reader()
			if err != nil {
				glog.Infoln(err)

				f.poller.Stop(desc)
				f.Driver.OnClose(driverConn)
				return
			}
			f.Driver.OnMessage(driverConn, msg)
		}()
	})
}
