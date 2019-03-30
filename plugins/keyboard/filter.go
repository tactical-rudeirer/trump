package keyboard

import (
	"trump/pkg/middleware"
	"trump/pkg/capture/pcap"
	"github.com/google/gopacket/layers"
)

var Filter = middleware.Plugin{
	Process:  FilterMsg,
	Priority: 0,
	Init:     InitFilter,
	Shutdown: func() {},
}

func InitFilter() {}

func FilterMsg(arg interface{}) interface{} {
	msg := arg.(pcap.USBData)
	if msg.EventType == layers.USBEventTypeComplete {
		return arg
	}
	return nil
}

