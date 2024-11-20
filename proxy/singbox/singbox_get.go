package singbox

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

func (s SBDetours) setData(config *users.Config) {

	for i := range s.Detours {
		for j := range config.Inbounds {
			if s.Detours[i].Detour != config.Inbounds[j].Tag {
				continue
			}

			if config.Inbounds[j].Protocol == "shadowsocks" {
				if len(config.Inbounds[j].Users) != 1 {
					continue
				}

				jsonStr := `{` +
					`"type":"shadowsocks",` +
					`"method":"` + config.Inbounds[j].Users[0].Method + `",` +
					`"password":"` + config.Inbounds[j].Users[0].Password + `"` +
					`}`

				config.Inbounds[j].Hide = true

				x := s.Detours[i].Index
				(*config.Inbounds[x].Shadowtls).DetourProxy = jsonStr
			}

		}
	}

}

func (inbound Inbound) getData(usersInbound *users.Inbound) string {
	protocol := inbound.Type
	usersInbound.Protocol = protocol

	if inbound.Listen == "::" {
		usersInbound.Port = inbound.Port
	}

	switch protocol {
	case "vmess", "vless", "trojan":
		usersInbound.Network = inbound.Transport.Type
		usersInbound.FixedSecurity = false

		if usersInbound.Network != "" {
			usersInbound.Transport = new(users.Transport)
		}

		switch usersInbound.Network {
		case "":
			usersInbound.Network = "tcp"
		case "grpc":
			(*usersInbound.Transport).ServiceName = inbound.Transport.ServiceName
		case "http":
			usersInbound.Transport.Path = inbound.Transport.Path
			if len(inbound.Transport.Host) > 0 {
				(*usersInbound.Transport).Host = ""
				for _, host := range inbound.Transport.Host {
					(*usersInbound.Transport).Host += host + ","
				}
				(*usersInbound.Transport).Host = strings.TrimRight((*usersInbound.Transport).Host, ",")
			}
		case "httpupgrade", "ws":
			(*usersInbound.Transport).Path = inbound.Transport.Path
		}

		if protocol == "vless" {
			if inbound.Tls.Reality.Enabled {
				usersInbound.Security = "reality"
				usersInbound.FixedSecurity = true
				usersInbound.Reality = new(users.Reality)
				(*usersInbound.Reality).Sni = inbound.Tls.ServerName
				if publicKey, err := change.GetPublicKey(inbound.Tls.Reality.PrivateKey); err == nil {
					(*usersInbound.Reality).PublicKey = publicKey
				}
				if len(inbound.Tls.Reality.ShortId) > 0 {
					(*usersInbound.Reality).ShortId = inbound.Tls.Reality.ShortId[0]
				}
			}
		}

		if inbound.Tls.Enabled && usersInbound.Security == "" {
			usersInbound.Security = "tls"
			usersInbound.FixedSecurity = true
			if len(inbound.Tls.Alpn) > 0 {
				usersInbound.Tls = new(users.Tls)
				(*usersInbound.Tls).Alpn = ""
				for _, alpn := range inbound.Tls.Alpn {
					(*usersInbound.Tls).Alpn += alpn + ","
				}
				(*usersInbound.Tls).Alpn = strings.TrimRight((*usersInbound.Tls).Alpn, ",")
			}
		}

		usersInbound.Fingerprint = setup.ConfigData.Users.UtlsFp

		return protocol
	case "hysteria2":
		usersInbound.Tls = new(users.Tls)
		(*usersInbound.Tls).Alpn = "h3"
		usersInbound.Security = "tls"
		return protocol
	case "tuic":
		usersInbound.Tls = new(users.Tls)
		(*usersInbound.Tls).Alpn = "h3"
		usersInbound.CongestionControl = inbound.CongestionControl
		usersInbound.Security = "tls"
		return protocol
	case "shadowsocks":
		if inbound.Password != "" {
			path, err := random.GenerateStrings(16)
			if err != nil {
				fmt.Println("随机url路径错误:", err)
				return ""
			}
			usersInbound.Users = append(usersInbound.Users, users.User{
				Name:     proxy.OnlyName + "-" + fmt.Sprintf("%d", inbound.Port),
				Password: inbound.Password,
				Method:   inbound.Method,
				Static:   false,
				UserPath: path,
			})
		}
		return protocol
	case "shadowtls":
		usersInbound.Shadowtls = new(users.Shadowtls)

		(*usersInbound.Shadowtls).Version = inbound.Version
		(*usersInbound.Shadowtls).Sni = inbound.Handshake.Server
		(*usersInbound.Shadowtls).Detour = inbound.Detour
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
	// var hides []string

	sbDetours := SBDetours{}

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

		total = len(config.Inbounds[i].Users)

		for j := range config.Inbounds[i].Users {

			if config.Inbounds[i].Users[j].Name == "" && total != 1 {
				continue
			}

			if config.Inbounds[i].Users[j].Name == "" {
				//name = proxy.OnlyName + "-" + fmt.Sprintf("%d", i)
				name = proxy.OnlyName + "-" + fmt.Sprintf("%d", config.Inbounds[i].Port)
			} else {
				name = config.Inbounds[i].Users[j].Name
			}

			path, err = random.GenerateStrings(16)
			if err != nil {
				fmt.Println("随机url路径错误:", err)
				return err
			}

			newUsersInbound.Users = append(newUsersInbound.Users, users.User{
				Name:     name,
				Static:   false,
				UserPath: path,
			})

			n := len(newUsersInbound.Users) - 1
			switch protocol {
			case "vmess":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Users[j].UUID
			case "vless":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Users[j].UUID
				newUsersInbound.Users[n].Flow = config.Inbounds[i].Users[j].Flow
			case "trojan", "shadowsocks", "shadowtls", "hysteria2":
				newUsersInbound.Users[n].Password = config.Inbounds[i].Users[j].Password
				if protocol == "shadowsocks" {
					newUsersInbound.Users[n].Method = config.Inbounds[i].Method
				}

			case "tuic":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Users[j].UUID
				newUsersInbound.Users[n].Password = config.Inbounds[i].Users[j].Password
			}
		}

		if protocol == "shadowtls" {
			sbDetours.Detours = append(sbDetours.Detours, Detour{
				Index:  len(usersConfig.Inbounds),
				Detour: newUsersInbound.Shadowtls.Detour,
			})
		}

		if len(newUsersInbound.Users) == 0 {
			newUsersInbound.Hide = true
		} else {
			newUsersInbound.Hide = false
		}

		usersConfig.Inbounds = append(usersConfig.Inbounds, newUsersInbound)
		newUsersInbound = users.Inbound{}
	}

	sbDetours.setData(&usersConfig)

	if setup.ConfigData.Static.Enabled {
		usersConfig.SetStaticUrl()
	}

	//fmt.Println(uupConfig)
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

		if len(config.Inbounds[i].Users) == 1 {
			if config.Inbounds[i].Users[0].Name == "" {
				switch config.Inbounds[i].Type {
				case "vmess":
					p.UserUUID = config.Inbounds[i].Users[0].UUID
				case "vless":
					p.UserUUID = config.Inbounds[i].Users[0].UUID
					p.UserFlow = config.Inbounds[i].Users[0].Flow
				case "trojan", "shadowsocks", "shadowtls", "hysteria2":
					p.UserPassword = config.Inbounds[i].Users[0].Password
				case "tuic":
					p.UserUUID = config.Inbounds[i].Users[0].UUID
					p.UserPassword = config.Inbounds[i].Users[0].Password
				}

				break OuterLoop
			}
		} else if config.Inbounds[i].Type == "shadowsocks" {
			if len(config.Inbounds[i].Users) == 0 && config.Inbounds[i].Password != "" {
				p.UserPassword = config.Inbounds[i].Password
				break OuterLoop
			}
		}

		for j := range config.Inbounds[i].Users {
			if config.Inbounds[i].Users[j].Name == userName {
				switch config.Inbounds[i].Type {
				case "vmess":
					p.UserUUID = config.Inbounds[i].Users[j].UUID
				case "vless":
					p.UserUUID = config.Inbounds[i].Users[j].UUID
					p.UserFlow = config.Inbounds[i].Users[j].Flow
				case "trojan", "shadowsocks", "shadowtls", "hysteria2":
					p.UserPassword = config.Inbounds[i].Users[j].Password

				case "tuic":
					p.UserUUID = config.Inbounds[i].Users[j].UUID
					p.UserPassword = config.Inbounds[i].Users[j].Password
				}
				break OuterLoop
			}
		}
	}
}
