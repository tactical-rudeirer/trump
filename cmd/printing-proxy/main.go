package main

import (
	"trump/pkg/proxy"
	"trump/plugins/printer"
	"trump/plugins/keyboard-filter"
)

func main() {
	proxy.RunProxy(0, 0, keyboard_filter.Plugin, printer.Plugin)
}
