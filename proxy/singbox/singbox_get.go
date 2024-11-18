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
					`"method":"` + config.Inbounds[j].Method + `",` +
					`"password":"` + config.Inbounds[j].Users[0].Password + `"` +
					`}`

				config.Inbounds[j].Hide = true

				x := s.Detours[i].Index
				config.Inbounds[x].Shadowtls = jsonStr
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
		switch inbound.Transport.Type {
		case "":
			usersInbound.Network = "tcp"
		case "grpc":
			usersInbound.ServiceName = inbound.Transport.ServiceName
		case "http":
			usersInbound.Path = inbound.Transport.Path
			if len(inbound.Transport.Host) > 0 {
				usersInbound.Host = ""
				for _, host := range inbound.Transport.Host {
					usersInbound.Host += host + ","
				}
				usersInbound.Host = strings.TrimRight(usersInbound.Host, ",")
			}
		case "httpupgrade", "ws":
			usersInbound.Path = inbound.Transport.Path
		}

		if protocol == "vless" {
			if inbound.Tls.Reality.Enabled {
				usersInbound.Security = "reality"
				usersInbound.FixedSecurity = true
				usersInbound.Sni = inbound.Tls.ServerName
				if publicKey, err := change.GetPublicKey(inbound.Tls.Reality.PrivateKey); err == nil {
					usersInbound.PublicKey = publicKey
				}
				if len(inbound.Tls.Reality.ShortId) > 0 {
					usersInbound.ShortId = inbound.Tls.Reality.ShortId[0]
				}
			}
		}

		if inbound.Tls.Enabled && usersInbound.Security == "" {
			usersInbound.Security = "tls"
			usersInbound.FixedSecurity = true
			if len(inbound.Tls.Alpn) > 0 {
				usersInbound.Alpn = ""
				for _, alpn := range inbound.Tls.Alpn {
					usersInbound.Alpn += alpn + ","
				}
				usersInbound.Alpn = strings.TrimRight(usersInbound.Alpn, ",")
			}
		}

		usersInbound.Fingerprint = setup.ConfigData.Users.UtlsFp

		return protocol
	case "hysteria2":
		usersInbound.Alpn = "h3"
		usersInbound.Security = "tls"
		return protocol
	case "tuic":
		usersInbound.Alpn = "h3"
		usersInbound.CongestionControl = inbound.CongestionControl
		usersInbound.Security = "tls"
		return protocol
	case "shadowsocks":
		usersInbound.Method = inbound.Method
	case "shadowtls":
		usersInbound.Version = inbound.Version
		usersInbound.Sni = inbound.Handshake.Server
		usersInbound.Detour = inbound.Detour
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

		if protocol == "shadowsocks" {

			if config.Inbounds[i].Password != "" {
				path, err = random.GenerateStrings(16)
				if err != nil {
					fmt.Println("随机url路径错误:", err)
					return err
				}
				newUsersInbound.Users = append(newUsersInbound.Users, users.User{
					Name:     proxy.OnlyName + "-" + fmt.Sprintf("%d", config.Inbounds[i].Port),
					Password: config.Inbounds[i].Password,
					Static:   false,
					UserPath: path,
				})
			}
		}

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
			case "tuic":
				newUsersInbound.Users[n].UUID = config.Inbounds[i].Users[j].UUID
				newUsersInbound.Users[n].Password = config.Inbounds[i].Users[j].Password
			}
		}
		// if len(newUsersInbound.Users) > 0 {
		// 	usersConfig.Inbounds = append(usersConfig.Inbounds, newUsersInbound)
		// }

		if protocol == "shadowtls" {
			sbDetours.Detours = append(sbDetours.Detours, Detour{
				Index:  len(usersConfig.Inbounds),
				Detour: newUsersInbound.Detour,
			})
			// hides = append(hides, newUsersInbound.Detour)
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
	// for i := range hides {
	// 	for j := range usersConfig.Inbounds {
	// 		if usersConfig.Inbounds[j].Tag == hides[i] {
	// 			usersConfig.Inbounds[j].Hide = true
	// 		}
	// 	}
	// }

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
				p.UserUUID = config.Inbounds[i].Users[0].UUID
				p.UserPassword = config.Inbounds[i].Users[0].Password
				p.UserFlow = config.Inbounds[i].Users[0].Flow
				break OuterLoop
			}
		}

		for j := range config.Inbounds[i].Users {
			if config.Inbounds[i].Users[j].Name == userName {
				p.UserUUID = config.Inbounds[i].Users[j].UUID
				p.UserPassword = config.Inbounds[i].Users[j].Password
				p.UserFlow = config.Inbounds[i].Users[j].Flow
				break OuterLoop
			}
		}
	}
}
