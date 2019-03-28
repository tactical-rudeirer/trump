package capture

import (
	"fmt"
	"time"
	"github.com/google/gopacket/pcap"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type USBData *layers.USB

var streams []chan USBData

func GetCaptureStream(buffer uint64) chan USBData {
	c := make(chan USBData, buffer)
	streams = append(streams, c)
	return c
}

func RunCapture(busID uint64, maxSize uint64, promisc bool, timeout time.Duration) func() {
	var handle *pcap.Handle
	inactive, err := pcap.NewInactiveHandle(fmt.Sprintf("usbmon%d", busID))
	if err != nil {
		log.Fatalf("could not create: %v", err)
	}
	if err = inactive.SetSnapLen(int(maxSize)); err != nil {
		log.Fatalf("could not set snap length: %v", err)
	} else if err = inactive.SetPromisc(promisc); err != nil {
		log.Fatalf("could not set promisc mode: %v", err)
	} else if err = inactive.SetTimeout(timeout); err != nil {
		log.Fatalf("could not set timeout: %v", err)
	}
	if handle, err = inactive.Activate(); err != nil {
		log.Fatal("PCAP Activate error: ", err)
	}
	src := gopacket.NewPacketSource(handle, gopacket.DecodersByLayerName["USB"])
	go func() {
		for p := range src.Packets() {
			usbLayer := p.Layer(layers.LayerTypeUSB).(*layers.USB)
			for _, s := range streams {
				s <- USBData(usbLayer)
			}
		}
	}()
	return func() {
		handle.Close()
		inactive.CleanUp()
	}
}
