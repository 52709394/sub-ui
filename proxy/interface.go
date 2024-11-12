package proxy

import "sub-ui/proxy/protocol"

var ConfigData Config
var OnlyName string

type Config interface {
	RenewData(string) error
	GetCurrentData(*protocol.Config, string, string)
}
