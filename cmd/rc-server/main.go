package main

import (
	"trump/pkg/proxy"
	"trump/plugins/printer"
	"trump/plugins/injector"
	"trump/plugins/keyboard-filter"
	"os"
	"strconv"
	"kegelbot/github.com/yanzay/log"
)

func main() {
	p, err := strconv.ParseInt(os.Args[1], 10, 16)
	if err != nil {
		log.Fatalf("failed to parse port number: %v", err)
	}
	injector.SERVER_PORT = uint16(p)
	injector.SERVER_CERT = os.Args[2]
	injector.SERVER_KEY = os.Args[3]
	injector.CA_CERT = os.Args[4]
	proxy.RunProxy(0, 0, keyboard_filter.Plugin, printer.Plugin, injector.Plugin)
}
