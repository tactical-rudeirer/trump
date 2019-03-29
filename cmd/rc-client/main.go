package main

import (
	"trump/pkg/capture"
	"github.com/google/gopacket/pcap"
	"net"
	"os"
	"log"
	"github.com/docker/libchan/spdy"
)

func main() {
	inS := capture.GetCaptureStream(100)
	defer capture.RunCapture(0, 1024, false, pcap.BlockForever)()
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	p, err := spdy.NewSpdyStreamProvider(conn, false)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	t := spdy.NewTransport(p)
	s, err := t.NewSendChannel()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	for msg := range inS {
		s.Send(msg.Payload)
	}
}
