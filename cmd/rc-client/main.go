package main

import (
	"trump/pkg/capture"
	"github.com/google/gopacket/pcap"
	"log"
	"github.com/docker/libchan/spdy"
	"github.com/google/gopacket/layers"
	"crypto/tls"
	"os"
)

func main() {
	inS := capture.GetCaptureStream(100)
	defer capture.RunCapture(0, 1024, false, pcap.BlockForever)()
	// Load server cert
	cert, err := tls.LoadX509KeyPair(os.Args[2], os.Args[3])
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	conn, err := tls.Dial("tcp", os.Args[1], &config)
	defer conn.Close()
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
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
		if msg.EventType == layers.USBEventTypeComplete {
			s.Send(msg.Payload)
		}
	}
}
