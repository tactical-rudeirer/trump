package middleware

import (
	"sort"
	"trump/pkg/inject"
	capture "trump/pkg/capture/pcap"
)

type Plugin struct {
	Process  func(interface{}) interface{}
	Priority int
	Init     func()
	Shutdown func()
}

var sortedPlugins []Plugin
var plugins map[int]Plugin

func RegisterPlugin(plugin Plugin) {
	plugins[plugin.Priority] = plugin
}

func ClearPlugins() {
	plugins = make(map[int]Plugin)
}

func InitMiddleware() func() {
	var keys []int
	for k := range plugins {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	sortedPlugins = []Plugin{}
	for _, k := range keys {
		sortedPlugins = append(sortedPlugins, plugins[k])
		plugins[k].Init()
	}
	return func() {
		for i := len(sortedPlugins)-1; i >= 0; i-- {
			sortedPlugins[i].Shutdown()
		}
	}
}

func ProcessMsg(msg interface{}, skipTo int) inject.Data {
	processedMsg := interface{}(msg)
	for _, p := range sortedPlugins {
		if p.Priority < skipTo {
			continue
		}
		processedMsg = p.Process(processedMsg)
		if processedMsg == nil {
			return inject.Data{}
		}
	}
	if processedMsg, ok := processedMsg.(capture.USBData); ok {
		return processedMsg.Payload
	}
	return processedMsg.(inject.Data)
}
