package setup

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sub-ui/change"
	"sub-ui/read"
)

var ConfigData Config

func GetData() {

	var file []byte
	var err error

	file, err = os.ReadFile("sub-ui.json")

	if err != nil {
		fmt.Println("错误:自定义配置文件无法读取!")
		os.Exit(1)
		return
	}

	err = json.Unmarshal(file, &ConfigData)

	if err != nil {
		fmt.Println("错误:自定义配置文件无法解析!")
		os.Exit(1)
		return
	}

	if read.CheckExistence(ConfigData.Proxy.Config) != "file" {
		fmt.Println("错误:proxy 配置不存在")
		os.Exit(1)
		return
	}

	var re *regexp.Regexp

	if ConfigData.Users.Domain != "" {

		re = regexp.MustCompile(`^([a-zA-Z0-9-]{1,63}\.){1,}[a-zA-Z]{2,}$`)

		if !re.MatchString(ConfigData.Users.Domain) {
			fmt.Println("错误:订阅域名不是有效的!")
			os.Exit(1)
			return
		}
	}

	re = regexp.MustCompile(`^xray|sing-box$`)

	if !re.MatchString(ConfigData.Proxy.Core) {
		fmt.Println("错误:不支持的内核配置")
		os.Exit(1)
		return
	}

	re = regexp.MustCompile(`^\?ed=\d{4}$`)

	if !re.MatchString(ConfigData.Users.Ws0Rtt) {
		ConfigData.Users.Ws0Rtt = ""
	}

	re = regexp.MustCompile(`^/[0-9a-zA-Z/]*[0-9a-zA-Z]*$`)

	if !re.MatchString(ConfigData.Server.Home.Url) {
		fmt.Println("错误:管理页面路径不是合法的!")
		os.Exit(1)
		return
	}

	if !re.MatchString(ConfigData.Server.Post.Renew) {
		fmt.Println("错误:更新用户url post 路径不是合法的!")
		os.Exit(1)
		return
	}

	if ConfigData.Backup.Enabled {
		if !re.MatchString(ConfigData.Server.Post.Backup) {
			fmt.Println("错误:更新用户url exclude 路径不是合法的!")
			os.Exit(1)
			return
		}
	}

	if !re.MatchString(ConfigData.Server.Post.Set) {
		fmt.Println("错误:proxy数据设置 post 路径不是合法的!")
		os.Exit(1)
		return
	}

	if !re.MatchString(ConfigData.Server.UserUrl) {
		fmt.Println("错误:用户前置 url 路径不是合法的!")
		os.Exit(1)
		return
	}

	if ConfigData.Download.Enabled {
		if !re.MatchString(ConfigData.Download.Url) {
			fmt.Println("错误:下载规则文件 url 路径不是合法的!")
			os.Exit(1)
			return
		}
	}

	re = regexp.MustCompile(`^([0-9]{1,5})$`)

	if re.MatchString(ConfigData.Server.Port) {

		portNumber, _ := strconv.Atoi(ConfigData.Server.Port)

		if portNumber < 0 || portNumber > 65535 {
			fmt.Println("错误:端口号超出范围(0-65535)")
			os.Exit(1)
			return
		}
	} else {
		fmt.Println("无效的端口格式")
		os.Exit(1)
		return
	}

	if ConfigData.Server.Home.User == "" || ConfigData.Server.Home.Password == "" {
		fmt.Println("错误:sui-ui,用户或密码不能为空的!")
		os.Exit(1)
		return
	}

	if ConfigData.Server.Cookie.Name == "" {
		CookieName = "session_token"
	} else {
		CookieName = url.QueryEscape(ConfigData.Server.Cookie.Name)
	}

	if ConfigData.Server.Cookie.Value == "" {
		CookieValue = url.QueryEscape(
			change.ToBase64(ConfigData.Server.Home.User + ConfigData.Server.Home.Password))
	} else {
		CookieValue = url.QueryEscape(ConfigData.Server.Cookie.Value)
	}

	if ConfigData.Server.Cookie.Day >= 1 && ConfigData.Server.Cookie.Day <= 90 {
		CookieDay = ConfigData.Server.Cookie.Day
	} else {
		CookieDay = 7
	}

	var isFp bool

	fp := []string{"chrome", "firefox", "safari", "ios", "android", "edge", "360", "qq", "random", "randomized"}

	for _, f := range fp {
		if ConfigData.Users.UtlsFp == f {
			isFp = true
			break
		}
	}

	if !isFp {
		fmt.Println("错误:不支持的Fingerprint!")
		os.Exit(1)
		return
	}

}

func SavedConfig() error {
	nowData, err := json.MarshalIndent(ConfigData, "", "  ")
	if err != nil {
		fmt.Println("文件:sub-ui.json")
		fmt.Println("JSON格式化错误:", err)
		return err
	}

	err = os.WriteFile("sub-ui.json", nowData, 0644)
	if err != nil {
		fmt.Println("文件:sub-ui.json")
		fmt.Println("文件写入错误:", err)
		return err
	}

	return nil
}
