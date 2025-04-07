package users

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"sub-ui/backup"
	"sub-ui/change"
	"sub-ui/proxy"
	"sub-ui/proxy/protocol"
	"sub-ui/setup"
)

func GetUrlData(proxyUrl string) (string, string) {

	re := regexp.MustCompile(`^(.*?)/([0-9a-zA-Z]*)/(.*?)\.((?:html|json))$`)
	if !re.MatchString(proxyUrl) {
		return "", ""
	}

	urlMatch := re.FindStringSubmatch(proxyUrl)

	var p protocol.Config
	userName := ""
	tag := ""

	for i := range ConfigData.Inbounds {

		if ConfigData.Inbounds[i].TagPath != urlMatch[1] {
			continue
		}
		urlName, _ := url.QueryUnescape(urlMatch[3])
		isUser := false

		for j := range ConfigData.Inbounds[i].Users {
			if ConfigData.Inbounds[i].Users[j].UserPath == urlMatch[2] && ConfigData.Inbounds[i].Users[j].Name == urlName {
				userName = ConfigData.Inbounds[i].Users[j].Name
				p.UserName = &ConfigData.Inbounds[i].Users[j].Name

				switch ConfigData.Inbounds[i].Protocol {
				case "vmess":
					p.UserUUID = &ConfigData.Inbounds[i].Users[j].UUID
				case "vless":
					p.UserUUID = &ConfigData.Inbounds[i].Users[j].UUID
					if ConfigData.Inbounds[i].Users[j].Flow != "" {
						p.UserFlow = &ConfigData.Inbounds[i].Users[j].Flow
					}
				case "trojan", "shadowsocks", "anytls", "shadowtls", "hysteria2":
					p.UserPassword = &ConfigData.Inbounds[i].Users[j].Password
				case "tuic":
					p.UserUUID = &ConfigData.Inbounds[i].Users[j].UUID
					p.UserPassword = &ConfigData.Inbounds[i].Users[j].Password
				}

				if ConfigData.Inbounds[i].Protocol == "shadowsocks" {
					p.Method = &ConfigData.Inbounds[i].Users[j].Method
				}

				isUser = true
				break
			}
		}

		if !isUser {
			fmt.Println("订阅警告:订阅连接不正确!")
			return "", ""
		}

		portNumber, _ := strconv.Atoi(ConfigData.Inbounds[i].Port)

		if ConfigData.Inbounds[i].Addr == "" ||
			portNumber < 1 ||
			portNumber > 65535 {
			fmt.Println("订阅警告:用户地址或端口未设置!")
			return "", ""
		}

		tag = ConfigData.Inbounds[i].Tag

		p.Addr = &ConfigData.Inbounds[i].Addr
		p.Port = &ConfigData.Inbounds[i].Port
		p.Protocol = &ConfigData.Inbounds[i].Protocol

		p.Fingerprint = &ConfigData.Inbounds[i].Fingerprint

		switch *p.Protocol {
		case "vmess", "vless", "trojan":

			p.Network = &ConfigData.Inbounds[i].Network

			switch *p.Network {
			case "grpc":
				if ConfigData.Inbounds[i].Transport != nil {
					if ConfigData.Inbounds[i].Transport.ServiceName != "" {
						p.ServiceName = &ConfigData.Inbounds[i].Transport.ServiceName
					}
				}

			case "http", "ws", "httpupgrade", "splithttp", "xhttp":
				if ConfigData.Inbounds[i].Transport != nil {
					if ConfigData.Inbounds[i].Transport.Host != "" {
						p.Host = &ConfigData.Inbounds[i].Transport.Host
					}
					if ConfigData.Inbounds[i].Transport.Path != "" {
						p.Path = &ConfigData.Inbounds[i].Transport.Path
					}
				}
			}

			p.Security = &ConfigData.Inbounds[i].Security

			if *p.Security == "reality" {
				if ConfigData.Inbounds[i].Reality != nil {
					p.Sni = &ConfigData.Inbounds[i].Reality.Sni
					p.PublicKey = &ConfigData.Inbounds[i].Reality.PublicKey
					p.ShortId = &ConfigData.Inbounds[i].Reality.ShortId
				}

			} else if *p.Security == "tls" {
				if ConfigData.Inbounds[i].Tls != nil {
					if ConfigData.Inbounds[i].Tls.Sni != "" {
						p.Sni = &ConfigData.Inbounds[i].Tls.Sni
					}

					if ConfigData.Inbounds[i].Tls.Alpn != "" {
						p.Alpn = &ConfigData.Inbounds[i].Tls.Alpn
					}
				}
			}

		case "anytls", "tuic", "hysteria2":
			if *p.Protocol == "tuic" {
				p.TuicCC = &ConfigData.Inbounds[i].CongestionControl
			}
			if ConfigData.Inbounds[i].Tls != nil {
				if ConfigData.Inbounds[i].Tls.Alpn != "" {
					p.Alpn = &ConfigData.Inbounds[i].Tls.Alpn
				}
			}
		case "shadowtls":

			p.Version = &ConfigData.Inbounds[i].Shadowtls.Version
			p.DetourProxy = &ConfigData.Inbounds[i].Shadowtls.DetourProxy
			p.Sni = &ConfigData.Inbounds[i].Shadowtls.Sni

		}

	}

	if setup.ConfigData.Proxy.RealTime {
		proxy.LConfigData.GetCurrentData(&p, tag, userName)
	}

	if p.UserUUID != nil && p.UserPassword != nil {
		fmt.Println("订阅警告:UUID或Password是无效的!")
		return "", ""
	}

	if urlMatch[4] == "json" {
		isUseBackup := false

		jsonStr := p.JsonUrl(setup.ConfigData.SingBox.MainTag)

		if jsonStr == "" {
			return "", "html"
		}

		if setup.ConfigData.Backup.Enabled {
			if backup.ProxySBData == "" || backup.SBSelectorOrUrlTestData == "" {
				goto jsonOut
			}

			for i := range setup.ConfigData.Backup.Excludes {

				if tag != setup.ConfigData.Backup.Excludes[i].Tag {
					continue
				}
				for _, n := range setup.ConfigData.Backup.Excludes[i].Users {
					if *p.UserName == n {
						goto jsonOut
					}
				}
			}

			jsonStr = backup.SBSelectorOrUrlTestData + jsonStr + ",\n" + backup.ProxySBData
			isUseBackup = true
		}
	jsonOut:
		urlSB := protocol.GenerateSBConfig(jsonStr, isUseBackup)
		if setup.ConfigData.SingBox.Format {
			return urlSB, "json"
		} else {
			return urlSB, "html"
		}

	}

	url := p.HttpUrl()

	if url == "" {
		goto urlOUt
	}

	if setup.ConfigData.Backup.Enabled {
		if backup.ProxyUrlData != "" {

			for i := range setup.ConfigData.Backup.Excludes {

				if tag != setup.ConfigData.Backup.Excludes[i].Tag {
					continue
				}
				for _, n := range setup.ConfigData.Backup.Excludes[i].Users {
					if *p.UserName == n {
						goto urlOUt
					}
				}
			}

			url += "\n" + backup.ProxyUrlData
		}
	}

urlOUt:
	urlBase64 := change.ToBase64(url)

	return urlBase64, "html"
}

func TagHttpString(inbound Inbound) (string, string, string) {
	tag := inbound.Tag
	strA := "tag:" + tag
	strB := "tag:" + tag

	switch inbound.Protocol {
	case "vmess", "vless", "trojan":
		strA += "(协议:" + inbound.Protocol +
			"+" + inbound.Network + "+"

		strB += "(协议:" + inbound.Protocol +
			"+" + inbound.Network + "+"

		if inbound.Security == "" {
			strA += "无传输层安全)"
			strB += "无传输层安全, "
		} else {
			strA += inbound.Security + ")"
			strB += inbound.Security + ", "
		}
	default:
		strA += "(协议:" + inbound.Protocol + ")"
		strB += "(协议:" + inbound.Protocol + " "
	}

	portNumber, _ := strconv.Atoi(inbound.Port)

	if portNumber < 1 || portNumber > 65535 ||
		inbound.Addr == "" {
		if inbound.Addr == "" {
			strB += "地址:,未设置, "
		} else {
			strB += "地址:" + inbound.Addr + ", "
		}

		if portNumber < 1 || portNumber > 65535 {
			strB += "端口:端口号超出范围, "
		} else {
			strB += "端口:" + inbound.Port + ", "
		}
		strB += "注意:地址或端口设置有误,无法生成订阅!)"
	} else {
		strB += "地址:" + inbound.Addr + ", " +
			"端口: " + inbound.Port + ")"
	}

	strC := ""

	if inbound.Addr != "" {
		strC += "地址: " + inbound.Addr + ", "
	} else {
		strC += "地址: 未设置, "
	}

	if inbound.Port != "" {
		strC += "端口: " + inbound.Port + ", "
	} else {
		strC += "端口: 未设置, "
	}

	if inbound.Security != "" {
		strC += "传输层安全: " + inbound.Security + ", "
	} else {
		strC += "传输层安全: 没, "
	}

	if inbound.Tls != nil {
		if inbound.Tls.Alpn == "" {
			strC += "alpn: " + inbound.Tls.Alpn
		}
	} else {
		strC += "alpn: 没 "
	}
	return strA, strB, strC
}

func UsersListHttp(subAddr string, setTagStr, usersLiSrt *string) {

	var urlpath, userName string

	//domain := "https://" + setup.ConfigData.Users.Domain

	if len(ConfigData.Inbounds) > 2 {
		*setTagStr += `
            <li>
                <p>所有tag</p>
            </li>
            <li>
                <label>地址:</label> <input id="addrInp-1" style="width: 150px;"> <label>端口:</label> <input id="portInp-1"
                    type="number" min="0" max="65535"> </br> </br>

                <label>传输层安全:</label>
                <select id="securitySel-1" style="height: 21px; width: 115px;">
                    <option value=""></option>
                    <option value="tls">tls</option>
                </select>

                <label>alpn:</label>
                <select id="alpnSel-1" style="height: 21px; width: 115px;">
                    <option value=""></option>
                    <option value="h3">h3</option>
                    <option value="h2">h2</option>
                    <option value="http/1.1">http/1.1</option>
                    <option value="h3,h2">h3,h2</option>
                    <option value="h2,http/1.1">h2,http/1.1</option>
                    <option value="h3,h2,http/1.1">h3,h2,http/1.1</option>
                </select>
                </br> </br>
                <button
				    onclick="setProxyData(-1,'所有tag')">保存设定</button>
            </li>
`
	}

	for i := range ConfigData.Inbounds {
		if ConfigData.Inbounds[i].Hide {
			continue
		}

		tag := ConfigData.Inbounds[i].Tag
		tagAID := fmt.Sprintf("tagA%d", i)
		tagBID := fmt.Sprintf("tagB%d", i)
		tagInfoID := fmt.Sprintf("tagInfoSpan%d", i)
		addrID := fmt.Sprintf("addrInp%d", i)
		portID := fmt.Sprintf("portInp%d", i)
		securityID := fmt.Sprintf("securitySel%d", i)
		alpnID := fmt.Sprintf("alpnSel%d", i)

		tagSrtA, tagSrtB, tagInfoSrt := TagHttpString(ConfigData.Inbounds[i])

		*setTagStr += `
            <li>
                <p id="` + tagAID + `">` + tagSrtA + `</p>
                <span id="` + tagInfoID + `">` + tagInfoSrt + `</span>
            </li>
            <li>
                <label>地址:</label> <input id="` + addrID + `" style="width: 150px;"> <label>端口:</label> <input id="` + portID + `"
                    type="number" min="0" max="65535"> </br> </br>

                <label>传输层安全:</label>
                <select id="` + securityID + `" style="height: 21px; width: 115px;">
                    <option value=""></option>
                    <option value="tls">tls</option>
                </select>

                <label>alpn:</label>
                <select id="` + alpnID + `" style="height: 21px; width: 115px;">
                    <option value=""></option>
                    <option value="h3">h3</option>
                    <option value="h2">h2</option>
                    <option value="http/1.1">http/1.1</option>
                    <option value="h3,h2">h3,h2</option>
                    <option value="h2,http/1.1">h2,http/1.1</option>
                    <option value="h3,h2,http/1.1">h3,h2,http/1.1</option>
                </select>
                </br> </br>
                <button
				    onclick="setProxyData(` + fmt.Sprintf("%d", i) + `,'` + tag + `')">保存设定</button>
            </li>
`

		*usersLiSrt += `<li><p id="` + tagBID + `">` + tagSrtB + `</p></li>`
		for j := range ConfigData.Inbounds[i].Users {
			urlpath = subAddr + setup.ConfigData.Server.UserUrl + "/" + ConfigData.Inbounds[i].TagPath + "/" +
				ConfigData.Inbounds[i].Users[j].UserPath + "/" + url.QueryEscape(ConfigData.Inbounds[i].Users[j].Name)
			userName = ConfigData.Inbounds[i].Users[j].Name
			strJson := fmt.Sprintf(`{"x":%d,"y":%d,"name":"%s"}`, i, j, userName)
			boxID := fmt.Sprintf("usersID%dU%d", i, j)
			*usersLiSrt += `
		<li><input type="checkbox" data-user-id='` + strJson + `' id="` + boxID + `"><label for="` + boxID + `">` + userName + `</label>
		<button
			onclick="showQRCode('sb','` + urlpath + `.json','` + userName + `')">sing-box</button>
		<button
			onclick="showQRCode('html','` + urlpath + `.html','` + userName + `')">htmlUrl</button>
		<button onclick="copyContent('` + urlpath + `.html')">htmlUrl</button>
	</li>
	`
		}
	}

}
