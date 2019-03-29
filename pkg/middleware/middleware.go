package middleware

import (
	"sort"
	"trump/pkg/inject"
	"trump/pkg/capture"
)

type Plugin struct {
	Process  func(interface{}) interface{}
	Priority int
	Init     func()
}
var sortedPlugins []Plugin
var plugins map[int]Plugin

func RegisterPlugin(plugin Plugin) {
	plugins[plugin.Priority] = plugin
}

func ClearPlugins() {
	plugins = make(map[int]Plugin)
}

func InitMiddleware() {
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
}

func ProcessMsg(msg interface{}, skipTo int) inject.Data {
	processedMsg := interface{}(msg)
	for _, p := range sortedPlugins {
		if p.Priority < skipTo {
			continue
		}
		processedMsg = p.Process(processedMsg)
	}
	if processedMsg, ok := processedMsg.(capture.USBData); ok {
		return processedMsg.Payload
	}
	return processedMsg.(inject.Data)
}
