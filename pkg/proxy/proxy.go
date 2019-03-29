package proxy

import (
	"trump/pkg/capture"
	"trump/pkg/inject"
	"github.com/google/gopacket/pcap"
	"trump/pkg/middleware"
)

func RunProxy(inID uint64, outID uint64, plugins ...middleware.Plugin) {
	inS := capture.GetCaptureStream(100)
	outS := inject.GetInjectionStream(100)
	defer capture.RunCapture(inID, 1024, false, pcap.BlockForever)()
	defer inject.RunInjector(outID)()
	middleware.ClearPlugins()
	for _, p := range plugins {
		middleware.RegisterPlugin(p)
	}
	middleware.InitMiddleware()
	for msg := range inS {
		outS <- middleware.ProcessMsg(msg)
	}
}
