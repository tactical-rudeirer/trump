package keyboard_filter

import (
	"trump/pkg/middleware"
	"trump/pkg/capture"
	"github.com/google/gopacket/layers"
)

var Plugin = middleware.Plugin{
	Process: ProcessMsg,
	Priority: 0,
	Init: Init,
	Shutdown: func() {},
}

func Init() {}

func ProcessMsg(arg interface{}) interface{} {
	msg := arg.(capture.USBData)
	if msg.EventType == layers.USBEventTypeComplete {
		return arg
	}
	return nil
}

