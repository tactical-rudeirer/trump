package inject

import (
	"os"
	"fmt"
	"log"
)

type Data []byte
var inStream chan Data

func GetInjectionStream(buf uint64) chan Data {
	inStream = make(chan Data, buf)
	return inStream
}

func RunInjector(id uint64) func () {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/hidg%d", id), os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("could not open hid device: %v", err)
	}
	go func() {
		for msg := range inStream {
			_, err := fd.Write(msg)
			log.Printf("failed to write to hid device: %v", err)
		}
	}()
	return func() {
		fd.Close()
	}
}
