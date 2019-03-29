package main

import (
	"trump/pkg/proxy"
	"trump/plugins/printer"
	"trump/plugins/injector"
)

func main() {
	proxy.RunProxy(0, 0, printer.Plugin, injector.Plugin)
}
