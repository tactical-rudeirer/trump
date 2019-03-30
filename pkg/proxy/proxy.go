package proxy

import (
	capture "trump/pkg/capture/pcap"
	"trump/pkg/inject"
	"github.com/google/gopacket/pcap"
	"trump/pkg/middleware"
)

type proxyMsg struct {
	msg interface{}
	skipTo int
}

var proxyChan chan proxyMsg

func RunProxy(inID uint64, outID uint64, plugins ...middleware.Plugin) {
	inS := capture.GetCaptureStream(100)
	proxyChan = make(chan proxyMsg, 100)
	outS := inject.GetInjectionStream(100)
	defer capture.RunCapture(inID, 1024, false, pcap.BlockForever)()
	defer inject.RunInjector(outID)()
	middleware.ClearPlugins()
	for _, p := range plugins {
		middleware.RegisterPlugin(p)
	}
	defer middleware.InitMiddleware()()
	go func() {
		for msg := range inS {
			proxyChan <- proxyMsg{msg, -1}
		}
	}()
	for msg := range proxyChan {
		outS <- middleware.ProcessMsg(msg.msg, msg.skipTo)
	}
}

func InjectMsg(msg interface{}, skipTo int) {
	proxyChan <- proxyMsg{msg, skipTo}
}
