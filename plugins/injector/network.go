package injector

import (
	"trump/pkg/middleware"
	"net"
	"log"
	"time"
	"github.com/docker/libchan/spdy"
	"trump/pkg/inject"
	"fmt"
	"trump/pkg/proxy"
)

var PLUGIN_PRIORITY = 10
var Plugin = middleware.Plugin{
	Process: ProcessMsg,
	Priority: PLUGIN_PRIORITY,
	Init: Init,
}

func Init() {
	//TODO: TLS
	listener, err := net.Listen("tcp", "0.0.0.0:9998")
	if err != nil {
		log.Fatalf("failed to start listening for connections: %v", err)
	}
	go func() {
		for {
			con, err := listener.Accept()
			if err == nil {
				p, err := spdy.NewSpdyStreamProvider(con, true)
				if err == nil {
					t := spdy.NewTransport(p)
					r, err := t.WaitReceiveChannel()
					if err == nil {
						msg := &inject.Data{}
						for {
							err = r.Receive(msg)
							if err != nil {
								fmt.Printf("failed to receive message: %v", err)
								continue
							}
							proxy.InjectMsg(msg, PLUGIN_PRIORITY)
						}
						continue
					}
				}
			}
			log.Printf("failed to accept connection: %v", err)
			time.Sleep(time.Second)
		}
	}()
}

func ProcessMsg(arg interface{}) interface{} {
	return arg
}
