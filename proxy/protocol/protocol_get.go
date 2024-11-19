package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func GetSBString() {
	getString := func(name string) string {
		if file, err := os.ReadFile(name); err == nil {
			return string(file)
		}
		return ""
	}

	SBStringData.UrlTest = getString("sing-box/url_test.json")
	SBStringData.Selector = getString("sing-box/selector.json")
	SBStringData.VmessWsTls = getString("sing-box/vmess_ws_tls.json")
	SBStringData.VlessTcpReality = getString("sing-box/vless_tcp_reality.json")
	SBStringData.VlessHttpReality = getString("sing-box/vless_http_reality.json")
	SBStringData.VlessGrpcReality = getString("sing-box/vless_grpc_reality.json")
	SBStringData.VlessTcpTls = getString("sing-box/vless_tcp_tls.json")
	SBStringData.TrojanTcpTls = getString("sing-box/trojan_tcp_tls.json")
	SBStringData.Hysteria2 = getString("sing-box/hysteria2.json")
	SBStringData.Tuic = getString("sing-box/tuic.json")
	SBStringData.Shadowtls = getString("sing-box/shadowtls.json")
	SBStringData.ShadowtlsSS = getString("sing-box/Shadowtls_ss.json")
	SBStringData.Shadowsocks = getString("sing-box/shadowsocks.json")

}

func (p Config) setSBData(proxyStr string, tag string) string {

	var sbTemp *template.Template
	var err error

	sbTemp, err = template.New("json").Parse(proxyStr)
	if err != nil {
		fmt.Println("sing-box 协议模板不正确!")
		fmt.Println(err)
		return ""
	}

	var alpn, host string

	if p.Alpn != "" {
		if strings.Contains(p.Alpn, ",") {
			alpn = `"alpn" : [`

			for _, a := range strings.Split(p.Alpn, ",") {
				alpn += `"` + a + `",`
			}
			alpn += `]`
			alpn = strings.Replace(alpn, ",]", "],", -1)
		} else {
			alpn = `"alpn" : ["` + p.Alpn + `"],`
		}
	}

	if p.Host != "" {

		if strings.Contains(p.Host, ",") {
			host = ""

			for _, h := range strings.Split(p.Host, ",") {
				host += `"` + h + `",`
			}

			host = strings.TrimRight(host, ", ")
		} else {
			host = `"` + p.Host + `"`
		}
	}

	proxyLate := struct {
		Tag         string
		Addr        string
		Port        string
		TuicCC      string
		Method      string
		Network     string
		HttpHost    template.HTML
		Path        string
		ServiceName string
		Alpn        template.HTML
		Sni         string
		Version     string
		PublicKey   string
		ShortId     string
		Fingerprint string
		Name        string
		UUID        string
		Password    string
	}{
		Tag:         tag,
		Addr:        p.Addr,
		Port:        p.Port,
		TuicCC:      p.TuicCC,
		Network:     p.Network,
		Method:      p.Method,
		HttpHost:    template.HTML(host),
		Path:        p.Path,
		ServiceName: p.ServiceName,
		Alpn:        template.HTML(alpn),
		Sni:         p.Sni,
		Version:     p.Version,
		PublicKey:   p.PublicKey,
		ShortId:     p.ShortId,
		Fingerprint: p.Fingerprint,
		UUID:        p.UserUUID,
		Password:    p.UserPassword,
	}

	var configBuf bytes.Buffer

	err = sbTemp.Execute(&configBuf, proxyLate)
	if err != nil {
		fmt.Println("sing-box 协议模板不正确!")
		fmt.Println(err)
		return ""
	}

	jsonString := strings.TrimSpace(configBuf.String())

	if !strings.HasSuffix(jsonString, "}") {
		fmt.Println("sing-box 协议模板不正确!")
		return ""
	}

	var js json.RawMessage

	err = json.Unmarshal([]byte(jsonString), &js)
	if err != nil {
		fmt.Println("singbox json 解析错误:", err)
		return ""
	}

	return jsonString
}
