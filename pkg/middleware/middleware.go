package middleware

import (
	"sort"
	"trump/pkg/capture"
	"trump/pkg/inject"
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

func ProcessMsg(msg capture.USBData) inject.Data {
	processedMsg := interface{}(msg)
	for _, p := range sortedPlugins {
		processedMsg = p.Process(processedMsg)
	}
	return processedMsg.(inject.Data)
}
