package printer

import (
	capture "trump/pkg/capture/pcap"
	"time"
	"github.com/google/gopacket/layers"
	"fmt"
	"trump/pkg/hid"
	"trump/pkg/middleware"
)

var Plugin = middleware.Plugin{
	Process: ProcessMsg,
	Priority: 5,
	Init: Init,
	Shutdown: func() {},
}

var lastString = ""
var lastTime = time.Duration(-1)
var diff = time.Duration(-1)

func Init() {
	lastString = ""
	lastTime = time.Duration(-1)
	diff = time.Duration(-1)
}

func ProcessMsg(arg interface{}) interface{} {
	msg := arg.(capture.USBData)
	if msg.EventType == layers.USBEventTypeComplete {
		t := time.Duration(msg.TimestampSec)*time.Second + time.Duration(msg.TimestampUsec)*time.Microsecond
		if lastTime >= 0 {
			diff = t - lastTime
		}
		lastTime = t
		if diff >= 0 {
			for i := time.Duration(0); i < (diff-time.Duration(250)*time.Millisecond)*109/10/time.Second; i++ {
				fmt.Print(lastString)
			}
		}
		lastString = hid.UsbToString(msg)
		fmt.Print(lastString)
	}
	return arg
}
