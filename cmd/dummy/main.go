package main

import "fmt"
import (
	"trump/pkg/capture"
	"time"
	"github.com/google/gopacket/pcap"
	"trump/pkg/hid"
	"github.com/google/gopacket/layers"
)

func main() {
	s := capture.GetCaptureStream(100)
	defer capture.RunCapture(0, 1024, false, pcap.BlockForever)()
	go func() {
		lastString := ""
		lastTime := time.Duration(-1)
		diff := time.Duration(-1)
		for o := range s {
			if o.EventType == layers.USBEventTypeComplete {
				t := time.Duration(o.TimestampSec) * time.Second + time.Duration(o.TimestampUsec) * time.Microsecond
				if lastTime >= 0 {
					diff = t - lastTime
				}
				lastTime = t
				if diff >= 0 {
					for i := time.Duration(0); i < (diff - time.Duration(250) * time.Millisecond) * 109 / 10 / time.Second; i++ {
						fmt.Print(lastString)
					}
				}
				lastString = hid.UsbToString(o)
				fmt.Print(lastString)
			}
		}
	}()
	time.Sleep(10 * time.Minute)
}
