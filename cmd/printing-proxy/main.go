package main

import (
	"trump/pkg/proxy"
	"trump/plugins/printer"
)

func main() {
	proxy.RunProxy(0, 0, printer.Plugin)
}
