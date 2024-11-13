package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sub-ui/change"
	"sub-ui/download"
	"sub-ui/proxy/protocol"
	"sub-ui/setup"
	"time"
)

var ProxyUrlData string
var ProxySBData string
var SBSelectorOrUrlTestData string

func GetUrlTicker() {

	hour := setup.ConfigData.Backup.StartTime

	if hour < 0 || hour > 23 {
		hour = 1
	}
	day := setup.ConfigData.Backup.UpdateInterval
	if day == 0 {
		day = 1
	}

	download.WaitStart(hour)

	ticker := time.NewTicker(time.Duration(day) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			GetProxyUrl()
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

			nextEnd := next.Add(20 * time.Minute)
			if !now.After(next) || !now.Before(nextEnd) {
				go GetUrlTicker()
				return
			}
		}
	}

}

func GetProxyUrl() {
	var urlData string
	var sbData string
	var tags string
	var selectorOrUrlTest string

	for i := range setup.ConfigData.Backup.ProxyList {

		if setup.ConfigData.Backup.ProxyList[i].Url == "" ||
			setup.ConfigData.Backup.ProxyList[i].SBTag == "" {
			continue
		}

		url, json := getUrl(setup.ConfigData.Backup.ProxyList[i].Url, setup.ConfigData.Backup.ProxyList[i].SBTag)

		if url != "" {
			urlData += url + "\n"
		}

		if json != "" {
			sbData += json + ",\n"
			tags += `"` + setup.ConfigData.Backup.ProxyList[i].SBTag + `",`
		}

	}

	if urlData != "" {
		ProxyUrlData = strings.TrimSpace(urlData)
	}

	if sbData == "" {
		return
	}

	tags = strings.TrimSuffix(strings.TrimSpace(tags), ",")

	tags = `"` + setup.ConfigData.SingBox.MainTag + `",` + tags

	if setup.ConfigData.Backup.SBSelector {
		if data, err := setSBTags(protocol.SBStringData.Selector, tags); err == nil {
			selectorOrUrlTest += data + ",\n"
		}
	}

	if data, err := setSBTags(protocol.SBStringData.UrlTest, tags); err == nil {
		selectorOrUrlTest += data + ",\n"
	}

	SBSelectorOrUrlTestData = selectorOrUrlTest
	ProxySBData = strings.TrimSuffix(strings.TrimSpace(sbData), ",")
}

func setSBTags(data string, tags string) (string, error) {
	var sbTemp *template.Template
	var err error

	sbTemp, err = template.New("json").Parse(data)
	if err != nil {
		fmt.Println("sing-box 协议模板不正确!")
		fmt.Println(err)
		return "", err
	}

	tagsLate := struct {
		Tags template.HTML
	}{
		Tags: template.HTML(tags),
	}

	var dataBuf bytes.Buffer

	err = sbTemp.Execute(&dataBuf, tagsLate)
	if err != nil {
		fmt.Println("sing-box 协议模板不正确!")
		fmt.Println(err)
		return "", err
	}

	jsonStr := strings.TrimSpace(dataBuf.String())

	if !strings.HasSuffix(jsonStr, "}") {
		fmt.Println("sing-box 协议模板不正确!")
		return "", fmt.Errorf("sing-box 协议模板不正确!")
	}

	return jsonStr, nil

}

func getUrl(url string, tag string) (string, string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应内容失败:", err)
		return "", ""
	}

	return setSBData(string(body), tag)
}

func setSBData(_url string, tag string) (string, string) {
	var proxyUrl string
	var proxyData string
	var match []string
	var err error
	var p protocol.Config
	re := regexp.MustCompile(`^(.*)\://`)

	if !re.MatchString(_url) {

		proxyUrl, err = change.ToUnicode(_url)
		if err != nil {
			if !re.MatchString(proxyUrl) {
				return "", ""
			}
		}

	} else {
		proxyUrl = _url
	}

	if len(re.FindAllString(proxyUrl, -1)) != 1 {
		return proxyUrl, ""
	}

	match = re.FindStringSubmatch(proxyUrl)
	protocol := match[1]

	switch protocol {
	case "vmess", "vless", "trojan", "hysteria2":
		re := regexp.MustCompile(`^(?:vmess|vless|trojan|hysteria2)\://(.*?)@(.*?)\:(\d{1,5})?`)

		if !re.MatchString(proxyUrl) {
			if protocol == "vmess" {
				vmess := strings.TrimPrefix(proxyUrl, "vmess://")
				return proxyUrl, getVmessData(vmess, tag)
			}
			return proxyUrl, ""
		}
		match = re.FindStringSubmatch(proxyUrl)

		p.Protocol = protocol
		if protocol == "vless" || protocol == "vmess" {
			p.UserUUID = match[1]
		} else {
			p.UserPassword, _ = url.PathUnescape(match[1])
		}

		p.Addr = match[2]
		p.Port = match[3]

		proxyData = strings.TrimPrefix(proxyUrl, match[0])
		return proxyUrl, getProxyData(p, proxyData, tag)
	case "tuic":
		re := regexp.MustCompile(`^tuic\://(.*?)\:(.*?)@(.*?)\:(\d{1,5})?`)
		if !re.MatchString(proxyUrl) {
			return proxyUrl, ""
		}
		p.UserUUID = match[1]
		p.UserPassword, _ = url.PathUnescape(match[1])
		p.Addr = match[3]
		p.Port = match[4]

		proxyData = strings.TrimPrefix(proxyUrl, match[0])
		return proxyUrl, getProxyData(p, proxyData, tag)

	}

	return proxyUrl, ""
}

func getVmessData(strData string, tag string) string {
	var jsonStr string
	var vmessJson map[string]string
	var err error
	var p protocol.Config

	jsonStr, err = change.ToUnicode(strData)

	if err != nil {
		return ""
	}
	err = json.Unmarshal([]byte(jsonStr), &vmessJson)

	if err != nil {
		return ""
	}

	p.Protocol = "vmess"
	p.Addr = vmessJson["add"]
	p.Port = vmessJson["port"]
	p.UserUUID = vmessJson["id"]
	p.Network = vmessJson["net"]
	p.Host, _ = url.QueryUnescape(vmessJson["host"])
	p.Path, _ = url.QueryUnescape(vmessJson["path"])
	p.Security = vmessJson["tls"]
	p.Sni = vmessJson["sni"]
	p.Alpn, _ = url.QueryUnescape(vmessJson["alpn"])
	p.Fingerprint = vmessJson["fp"]

	return p.JsonUrl(tag)
}

func getProxyData(p protocol.Config, proxyData string, tag string) string {

	var data string
	var err error

	parts := strings.Split(proxyData, "#")

	proxyData = parts[0]

	for _, v := range strings.Split(proxyData, "&") {

		if v == "" {
			continue
		}

		if data, err = getRegData(v, `type=(.*)`); err == nil {
			p.Network = data
		}

		if data, err = getRegData(v, `security=(.*)`); err == nil {
			p.Security = data
		}

		if data, err = getRegData(v, `alpn=(.*)`); err == nil {
			p.Alpn, _ = url.PathUnescape(data)
		}

		if data, err = getRegData(v, `host=(.*)`); err == nil {
			p.Host, _ = url.PathUnescape(data)
		}

		if data, err = getRegData(v, `path=(.*)`); err == nil {
			p.Path, _ = url.PathUnescape(data)
		}

		if data, err = getRegData(v, `sni=(.*)`); err == nil {
			p.Sni = data
		}

		if data, err = getRegData(v, `fp=(.*)`); err == nil {
			p.Fingerprint = data
		}

		if data, err = getRegData(v, `pbk=(.*)`); err == nil {
			p.PublicKey, _ = url.PathUnescape(data)
		}

		if data, err = getRegData(v, `sid=(.*)`); err == nil {
			p.ShortId = data
		}

		if data, err = getRegData(v, `serviceName=(.*)`); err == nil {
			p.ServiceName, _ = url.PathUnescape(data)
		}

		if data, err = getRegData(v, `flow=(.*)`); err == nil {
			p.UserFlow = data
		}

		if data, err = getRegData(v, `congestion_control=(.*)`); err == nil {
			p.TuicCC = data
		}

	}

	return p.JsonUrl(tag)
}

func getRegData(data string, regStr string) (string, error) {
	var re *regexp.Regexp

	re = regexp.MustCompile(regStr)

	if !re.MatchString(data) {
		return "", fmt.Errorf("数据不存在")
	}
	match := re.FindStringSubmatch(data)
	return match[1], nil
}
