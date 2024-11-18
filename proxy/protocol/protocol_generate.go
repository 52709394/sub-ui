package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"sub-ui/change"
	"sub-ui/setup"
)

func (p Config) HttpUrl() string {
	var httpStr string
	var password string
	var path, host, serviceName string
	var publicKey, alpn string
	var name string
	var vless, vmess string

	if p.Protocol == "vmess" && setup.ConfigData.Users.VmessModel == "old" {

		if (p.Network == "ws" || p.Network == "httpupgrade") &&
			setup.ConfigData.Proxy.Core == "xray" {
			p.Path += setup.ConfigData.Users.Ws0Rtt
		}

		return p.vmessUrl()
	}

	proxyMod := p.Protocol

	switch p.Protocol {
	case "vmess", "vless", "trojan":
		proxyMod += "+"
		proxyMod += p.Network + "+"
		proxyMod += p.Security
	}

	// if p.Protocol == "vless" || p.Protocol == "trojan" {
	// 	proxyMod += "+"
	// 	proxyMod += p.Network + "+"
	// 	proxyMod += p.Security
	// }

	proxyMod = strings.ToLower(proxyMod)

	name = url.QueryEscape(p.UserName)

	if p.Alpn != "" {
		alpn = "&alpn=" + url.QueryEscape(p.Alpn)
	}

	if p.Security == "reality" {
		publicKey = url.QueryEscape(p.PublicKey)
	}

	if p.Host != "" {
		host = "&host=" + url.QueryEscape(p.Host)
	}

	if p.Path != "" {
		if strings.HasPrefix(p.Path, "/") {
			path = p.Path
		} else {
			path = "/" + p.Path
		}
		if (p.Network == "ws" || p.Network == "httpupgrade") &&
			setup.ConfigData.Proxy.Core == "xray" {
			p.Path += setup.ConfigData.Users.Ws0Rtt
		}
		path = "&path=" + url.QueryEscape(path)
	}

	if p.ServiceName != "" {
		serviceName = "&mode=gun&serviceName=" + url.QueryEscape(p.ServiceName)
	}

	if p.UserPassword != "" {
		password = url.QueryEscape(p.UserPassword)
	}

	vmess = `vmess://` + p.UserUUID + `@` + p.Addr + `:` + p.Port + `?encryption=none`

	if proxyMod == "vmess+ws+tls" || proxyMod == "vmess+httpupgrade+tls" {

		httpStr = vmess + `&security=tls` + alpn + `&fp=` + p.Fingerprint + `&type=` + p.Network + path + `#` + name
		return httpStr
	}

	vless = `vless://` + p.UserUUID + `@` + p.Addr + `:` + p.Port + `?encryption=none`

	if proxyMod == "vless+tcp+reality" {
		if p.UserFlow != "xtls-rprx-vision" {
			fmt.Println("订阅警告:xtls-rprx-reality,没开启xtls-rprx-vision!")
			return ""
		}

		httpStr := vless + `&flow=xtls-rprx-vision&security=reality&sni=` + p.Sni + `&fp=` + p.Fingerprint + `&pbk=` + publicKey + `&sid=` + p.ShortId + `&type=tcp&headerType=none#` + name
		return httpStr
	}

	if proxyMod == "vless+xhttp+reality" {

		httpStr = vless + `&security=reality&sni=` + p.Sni + `&fp=` + p.Fingerprint + `&pbk=` + publicKey + `&sid=` + p.ShortId + `&type=xhttp` + path + host + `#` + name
		return httpStr
	}

	if proxyMod == "vless+http+reality" {

		httpStr = vless + `&security=reality&sni=` + p.Sni + `&fp=` + p.Fingerprint + `&pbk=` + publicKey + `&sid=` + p.ShortId + `&type=http` + host + `#` + name
		return httpStr
	}

	if proxyMod == "vless+grpc+reality" {

		httpStr = vless + `&security=reality&sni=` + p.Sni + `&fp=` + p.Fingerprint + `&pbk=` + publicKey + `&sid=` + p.ShortId + `&type=grpc` + serviceName + `#` + name

		return httpStr
	}

	if proxyMod == "vless+tcp+tls" {

		if p.UserFlow != "xtls-rprx-vision" {
			fmt.Println("订阅警告:vless+tcp+tls,没开启xtls-rprx-vision!")
			return ""
		}
		httpStr = vless + `&flow=xtls-rprx-vision&security=tls` + alpn + `&fp=` + p.Fingerprint + `&type=tcp&headerType=none#` + name

		return httpStr
	}

	if proxyMod == "vless+xhttp+tls" {

		httpStr = vless + `&security=tls` + alpn + `&fp=` + p.Fingerprint + `&type=xhttp` + path + host + `#` + name

		return httpStr
	}

	switch proxyMod {
	case "vless+ws+tls", "vless+httpupgrade+tls", "vless+splithttp+tls":
		httpStr = vless + `&security=tls` + alpn + `&fp=` + p.Fingerprint + `&type=` + p.Network + path + `#` + name

		return httpStr

	}

	if proxyMod == "trojan+tcp+tls" {
		httpStr = `trojan://` + password + `@` + p.Addr + `:` + p.Port + `?security=tls` + alpn + `&fp=` + p.Fingerprint + `&type=tcp&headerType=none#` + name
		return httpStr

	}

	if proxyMod == "hysteria2" {
		httpStr = `hysteria2://` + password + `@` + p.Addr + `:` + p.Port + `/?alpn=h3&insecure=0#` + name
		return httpStr
	}

	if proxyMod == "tuic" {
		httpStr = `tuic://` + p.UserUUID + `:` + password + `@` + p.Addr + `:` + p.Port + `?alpn=h3&congestion_control=` + p.TuicCC + `#` + name
		return httpStr
	}

	fmt.Println("订阅警告:Url协议暂不支持!")
	return ""

}

func getDetourData(str, tag string) string {

	var jsonMap map[string]string
	var p Config

	err := json.Unmarshal([]byte(str), &jsonMap)

	if err != nil {
		return ""
	}

	if jsonMap["type"] == "shadowsocks" {
		p.Protocol = "shadowtls_ss"
		p.Method = jsonMap["method"]
		p.UserPassword = jsonMap["password"]
		return p.JsonUrl(tag)
	}

	return ""
}

func (p Config) JsonUrl(tag string) string {

	proxyMod := p.Protocol

	protocols := []string{"vmess", "vless", "trojan"}

	for _, protocol := range protocols {
		if proxyMod == protocol {
			proxyMod += "+"
			proxyMod += p.Network + "+"
			proxyMod += p.Security
			break
		}
	}

	proxyMod = strings.ToLower(proxyMod)

	if proxyMod == "vmess+tcp+tls" {

	}

	if proxyMod == "vmess+ws+tls" || proxyMod == "vmess+httpupgrade+tls" {
		return p.setSBData(SBStringData.VmessWsTls, tag)
	}

	if proxyMod == "vless+tcp+reality" {
		if p.UserFlow != "xtls-rprx-vision" {
			fmt.Println("sing-box订阅警告:xtls-rprx-reality,没开启xtls-rprx-vision!")
			return ""
		}
		return p.setSBData(SBStringData.VlessTcpReality, tag)
	}

	if proxyMod == "vless+http+reality" {
		return p.setSBData(SBStringData.VlessHttpReality, tag)
	}

	if proxyMod == "vless+grpc+reality" {
		return p.setSBData(SBStringData.VlessGrpcReality, tag)
	}

	if proxyMod == "vless+tcp+tls" {

		if p.UserFlow != "xtls-rprx-vision" {
			fmt.Println("sing-box订阅警告:vless+tcp+tls,没开启xtls-rprx-vision!")
			return ""
		}
		return p.setSBData(SBStringData.VlessTcpTls, tag)

	}

	if proxyMod == "vless+ws+tls" || proxyMod == "vless+httpupgrade+tls" {

	}

	if proxyMod == "trojan+tcp+tls" {
		return p.setSBData(SBStringData.TrojanTcpTls, tag)
	}

	if proxyMod == "hysteria2" {
		return p.setSBData(SBStringData.Hysteria2, tag)
	}

	if proxyMod == "tuic" {
		return p.setSBData(SBStringData.Tuic, tag)
	}

	if proxyMod == "shadowtls" {
		shadowtlsStr := getDetourData(p.Shadowtls, tag)

		if shadowtlsStr == "" {
			return ""
		}
		shadowtlsStr += ",\n" + p.setSBData(SBStringData.Shadowtls, tag)

		return shadowtlsStr
	}

	if proxyMod == "shadowtls_ss" {
		return p.setSBData(SBStringData.ShadowtlsSS, tag)
	}

	fmt.Println("sing-box订阅警告:协议暂不支持!")
	return ""
}

func (p Config) vmessUrl() string {
	var path string

	if p.Path != "" {
		if strings.HasPrefix(p.Path, "/") {
			path = p.Path
		} else {
			path = "/" + p.Path
		}
	}

	vmess := `
	{
    "v": "2",
    "ps": "` + p.UserName + `",
    "add": "` + p.Addr + `",
    "port": "` + p.Port + `",
    "id": "` + p.UserUUID + `",
    "aid": "0",
    "scy": "auto",
    "net": "` + p.Network + `",
    "type": "none",
    "host": "` + p.Host + `",
    "path": "` + path + `",
    "tls": "` + p.Security + `",
    "sni": "` + p.Sni + `",
    "alpn": "` + p.Alpn + `",
    "fp": "` + p.Fingerprint + `"
    }`

	var formatBuf bytes.Buffer

	if err := json.Indent(&formatBuf, []byte(vmess), "", "    "); err != nil {
		fmt.Println("生成 vmess url json 错误:", err)
		return ""
	}

	url := "vmess://" + change.ToBase64(formatBuf.String())

	return url

}

func GenerateSBConfig(proxysString string, isUseBackup bool) string {

	var err error
	var sbTemp *template.Template
	var js json.RawMessage

	sbTemp, err = template.ParseFiles(setup.ConfigData.SingBox.Config)
	if err != nil {
		fmt.Println(err)
	}

	var configBuf bytes.Buffer
	var outbound, detour, uiDetour string
	if isUseBackup {
		outbound = setup.ConfigData.Backup.SingBox.Outbound

		if outbound == "" {
			outbound = setup.ConfigData.SingBox.MainTag
		}

		detour = setup.ConfigData.Backup.SingBox.DownloadDetour

		if detour == "" {
			detour = setup.ConfigData.SingBox.MainTag
		}

		uiDetour = setup.ConfigData.Backup.SingBox.ExternalUiDownloadDetour

		if uiDetour == "" {
			uiDetour = setup.ConfigData.SingBox.MainTag
		}

	} else {
		outbound = setup.ConfigData.SingBox.MainTag
		detour = setup.ConfigData.SingBox.MainTag
		uiDetour = setup.ConfigData.SingBox.MainTag
	}

	proxyLate := struct {
		Proxys   template.HTML
		Outbound string
		Detour   string
		UiDetour string
	}{
		Proxys:   template.HTML(proxysString),
		Outbound: outbound,
		Detour:   detour,
		UiDetour: uiDetour,
	}

	err = sbTemp.Execute(&configBuf, proxyLate)
	if err != nil {
		fmt.Println(err)
	}

	configBytes := configBuf.Bytes()

	err = json.Unmarshal(configBytes, &js)
	if err != nil {
		fmt.Println("singbox json 解析错误:", err)
		return ""
	}

	var formatBuf bytes.Buffer

	if err = json.Indent(&formatBuf, configBytes, "", "    "); err != nil {
		fmt.Println("singbox json 格式化错误:", err)
		return ""
	}

	return formatBuf.String()
}
