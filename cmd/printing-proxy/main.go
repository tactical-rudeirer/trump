package main

import (
	"trump/pkg/proxy"
	"trump/plugins/printer"
	"trump/plugins/keyboard"
)

func main() {
	proxy.RunProxy(0, 0, keyboard.Filter, printer.Plugin)
}
