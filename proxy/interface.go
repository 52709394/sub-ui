package proxy

import "sub-ui/proxy/protocol"

var ConfigData Config
var LConfigData LConfig
var OnlyName string

type Config interface {
	RenewData(string) error
}

type LConfig interface {
	GetCurrentData(*protocol.Config, string, string)
}
