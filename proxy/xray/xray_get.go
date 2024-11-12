package xray

import (
	"fmt"
	"strings"
	"sub-ui/change"
	"sub-ui/proxy"
	"sub-ui/proxy/protocol"
	"sub-ui/random"
	"sub-ui/read"
	"sub-ui/setup"
	"sub-ui/users"
)

func (inbound Inbound) getData(usersInbound *users.Inbound) string {
	protocol := inbound.Protocol
	usersInbound.Protocol = protocol

	if inbound.Listen == "" || inbound.Listen == "0.0.0.0" {
		usersInbound.Port = inbound.Port
	}

	switch protocol {
	case "vmess", "vless", "trojan":

		usersInbound.Network = inbound.StreamSettings.Network
		usersInbound.FixedSecurity = false

		switch usersInbound.Network {
		case "raw":
			usersInbound.Network = "tcp"
		case "grpc":
			usersInbound.ServiceName = inbound.StreamSettings.GrpcSettings.ServiceName
		case "h2", "http":
			usersInbound.Network = "http"
			usersInbound.Path = inbound.StreamSettings.HttpSettings.Path
			if len(inbound.StreamSettings.HttpSettings.Host) > 0 {
				usersInbound.Host = ""
				for _, host := range inbound.StreamSettings.HttpSettings.Host {
					usersInbound.Host += host + ","
				}
				usersInbound.Host = strings.TrimRight(usersInbound.Host, ",")
			}
		case "ws":
			usersInbound.Path = inbound.StreamSettings.WsSettings.Path
		case "splithttp":
			usersInbound.Path = inbound.StreamSettings.SplithttpSettings.Path
		case "xhttp":
			usersInbound.Path = inbound.StreamSettings.XhttpSettings.Path
		case "httpupgrade":
			usersInbound.Path = inbound.StreamSettings.HttpupgradeSettings.Path
		}

		if protocol == "vless" {
			if inbound.StreamSettings.Security == "reality" {
				usersInbound.Security = "reality"
				usersInbound.FixedSecurity = true
				usersInbound.Sni = inbound.StreamSettings.RealitySettings.ServerNames[0]
				if publicKey, err := change.GetPublicKey(inbound.StreamSettings.RealitySettings.PrivateKey); err == nil {
					usersInbound.PublicKey = publicKey
				}
				if len(inbound.StreamSettings.RealitySettings.ShortIds) > 0 {
					usersInbound.ShortId = inbound.StreamSettings.RealitySettings.ShortIds[0]
				}
			}
		}

		if inbound.StreamSettings.Security == "tls" {
			usersInbound.Security = "tls"
			usersInbound.FixedSecurity = true
			if len(inbound.StreamSettings.TlsSettings.Alpn) > 0 {
				usersInbound.Alpn = ""
				for _, alpn := range inbound.StreamSettings.TlsSettings.Alpn {
					usersInbound.Alpn += alpn + ","
				}
				usersInbound.Alpn = strings.TrimRight(usersInbound.Alpn, ",")
			}
		}

		usersInbound.Fingerprint = setup.ConfigData.Users.UtlsFp

		return protocol
	}

	return ""
}

func (config Config) RenewData(mod string) error {

	usersConfig := users.Config{}

	read.GetJsonData(setup.ConfigData.Proxy.Config, &config)

	var newUsersInbound users.Inbound
	var path string
	var base64 string
	var total int
	var name string
	var err error

	path = ""

	for i := range config.Inbounds {

		if config.Inbounds[i].Tag == "" {
			continue
		}

		protocol := config.Inbounds[i].getData(&newUsersInbound)

		if protocol == "" {
			newUsersInbound = users.Inbound{}
			continue
		}

		base64 = change.ToBase64(config.Inbounds[i].Tag)

		base64 = strings.ReplaceAll(base64, "+", "252B")
		base64 = strings.ReplaceAll(base64, "/", "252F")
		base64 = strings.ReplaceAll(base64, "=", "253D")

		newUsersInbound.Tag = config.Inbounds[i].Tag
		newUsersInbound.TagPath = base64

		total = len(config.Inbounds[i].Settings.Clients)

		for j := range config.Inbounds[i].Settings.Clients {

			if config.Inbounds[i].Settings.Clients[j].Email == "" && total != 1 {
				continue
			}

			if config.Inbounds[i].Settings.Clients[j].Email == "" {
				name = proxy.OnlyName + "-" + fmt.Sprintf("%d", i)
			} else {
				name = config.Inbounds[i].Settings.Clients[j].Email
			}

			path, err = random.GenerateStrings(16)
			if err != nil {
				fmt.Println("随机url路径错误:", err)
				return err
			}

			newUsersInbound.Users = append(newUsersInbound.Users, users.User{
				Name:     name,
				UserPath: path,
			})

			n := len(newUsersInbound.Users) - 1
			switch protocol {
			case "vmess":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Settings.Clients[j].Id
			case "vless":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Settings.Clients[j].Id
				newUsersInbound.Users[n].Flow = config.Inbounds[i].Settings.Clients[j].Flow
			case "trojan":
				newUsersInbound.Users[n].Password = config.Inbounds[i].Settings.Clients[j].Password

			}
		}
		// if len(newUsersInbound.Users) > 0 {
		// 	usersConfig.Inbounds = append(usersConfig.Inbounds, newUsersInbound)
		// }
		usersConfig.Inbounds = append(usersConfig.Inbounds, newUsersInbound)
		newUsersInbound = users.Inbound{}
	}

	path = ""
	if mod == "renew" {
		usersConfig.SetOldData()
	}

	err = usersConfig.SavedConfig()
	if err != nil {
		return err
	}

	return nil
}

func (config Config) GetCurrentData(p *protocol.Config, tag string, userName string) {
	read.GetJsonData(setup.ConfigData.Proxy.Config, &config)
OuterLoop:
	for i := range config.Inbounds {
		if config.Inbounds[i].Tag != tag {
			continue
		}

		if len(config.Inbounds[i].Settings.Clients) == 1 {
			if config.Inbounds[i].Settings.Clients[0].Email == "" {
				p.UserUUID = config.Inbounds[i].Settings.Clients[0].Id
				p.UserPassword = config.Inbounds[i].Settings.Clients[0].Password
				p.UserFlow = config.Inbounds[i].Settings.Clients[0].Flow
				break OuterLoop
			}
		}

		for j := range config.Inbounds[i].Settings.Clients {
			if config.Inbounds[i].Settings.Clients[j].Email == userName {
				p.UserUUID = config.Inbounds[i].Settings.Clients[j].Id
				p.UserPassword = config.Inbounds[i].Settings.Clients[j].Password
				p.UserFlow = config.Inbounds[i].Settings.Clients[j].Flow
				break OuterLoop
			}

		}
	}
}
