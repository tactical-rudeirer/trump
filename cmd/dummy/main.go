package main

import "fmt"
import (
	"trump/pkg/capture"
	"time"
	"github.com/google/gopacket/pcap"
)

func main() {
	s := capture.GetCaptureStream(100)
	go capture.RunCapture(0, 1024, false, pcap.BlockForever)
	go func() {
		for o := range s {
			fmt.Printf("%+v\n", o)
		}
	}()
	time.Sleep(10 * time.Second)
}
