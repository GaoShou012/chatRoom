package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/src/frontier"
	"github.com/src/frontier/frontier_ws"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	addr = flag.String("addr",":1234","")
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln, err := (&net.ListenConfig{KeepAlive: time.Minute}).Listen(ctx, "tcp", *addr)
	if err != nil {
		glog.Fatal(err)
	}

	driver := &frontier_ws.Driver{}

	fr,err := frontier.NewFrontier(ln,driver)
	if err != nil{
		panic(err)
	}
	fr.Start()
	fmt.Println("running")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch s := <-c; s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("got signal %s; stop server", s)
		case syscall.SIGHUP:
			glog.Infof("got signal %s; go to deamon", s)
			continue
		}
		if err := fr.Stop(); err != nil {
			glog.Errorf("stop server error: %v", err)
		}
		break
	}
}
