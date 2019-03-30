package dev_input

import (
	"time"
	"os"
	"encoding/binary"
	"fmt"
	"bytes"
	"log"
)

type DevInputData struct {
	Type  uint16
	Code  uint16
	Value int32
	Time  time.Time
}

var streams []chan DevInputData

func GetCaptureStream(buffer uint64) chan DevInputData {
	c := make(chan DevInputData, buffer)
	streams = append(streams, c)
	return c
}

func RunCapture(devID uint64) func() {
	f, err := os.Open(fmt.Sprintf("/dev/input/event%d", devID))
	if err != nil {
		log.Fatalf("failed to open /dev/input/event%d: %v", devID, err)
	}
	b := make([]byte, 24)
	go func() {
		for {
			f.Read(b)
			sec := binary.LittleEndian.Uint64(b[0:8])
			usec := binary.LittleEndian.Uint64(b[8:16])
			t := time.Unix(int64(sec), int64(usec))
			var value int32
			typ := binary.LittleEndian.Uint16(b[16:18])
			code := binary.LittleEndian.Uint16(b[18:20])
			binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)
			event := DevInputData{Type: typ, Code: code, Value: value, Time: t}
			for _, s := range streams {
				s <- event
			}
		}
	}()
	return func() {
		f.Close()
	}
}
