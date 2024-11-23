package serve

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sub-ui/backup"
	"sub-ui/download"
	"sub-ui/proxy"
	"sub-ui/proxy/protocol"
	"sub-ui/proxy/singbox"
	"sub-ui/proxy/xray"
	"sub-ui/random"
	"sub-ui/read"
	"sub-ui/setup"
	"sub-ui/users"
	"sync"
)

var Mu sync.Mutex
var ToggleContent string

func (s Server) setTag(w http.ResponseWriter, r *http.Request) {

	Mu.Lock()
	defer Mu.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, "无法读取请求!", http.StatusBadRequest)
		return
	}

	var response map[string]string
	var errStr string

	response, ToggleContent, errStr = users.SetTagData(body)

	if errStr != "" {
		jsonError(w, errStr, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s Server) renewUrl(w http.ResponseWriter, r *http.Request) {

	Mu.Lock()
	defer Mu.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, "无法读取请求!", http.StatusBadRequest)
		return
	}

	var renewUsers users.RenewUsers

	err = json.Unmarshal(body, &renewUsers)
	if err != nil {
		jsonError(w, "无法解析数据!", http.StatusBadRequest)
		return
	}

	if renewUsers.Mod == "reset" {

		err = renewUsers.SetUsersUrl()
		if err != nil {
			jsonError(w, "重置失败,重试或联系技术人员!", http.StatusBadRequest)
			return
		}

	} else if renewUsers.Mod == "static" {
		renewUsers.SetStaticUsers()

	} else {

		err = proxy.ConfigData.RenewData("renew")
		if err != nil {
			jsonError(w, "更新用户数据失败,重试或联系技术人员!", http.StatusBadRequest)
			return
		}

	}

	ToggleContent = "user"

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "successfully"}
	json.NewEncoder(w).Encode(response)
}

func (s Server) renewBackupUrl(w http.ResponseWriter, r *http.Request) {

	Mu.Lock()
	defer Mu.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, "无法读取请求!", http.StatusBadRequest)
		return
	}

	var info users.BackupInfo

	err = json.Unmarshal(body, &info)
	if err != nil {
		jsonError(w, "无法解析数据!", http.StatusBadRequest)
		return
	}

	if info.Mod == "exclude" {
		info.AddUsers()
		setup.SavedConfig()
	} else {
		backup.GetProxyUrl()
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "successfully"}
	json.NewEncoder(w).Encode(response)

}

func (s Server) usersUrl(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if !strings.HasPrefix(path, setup.ConfigData.Server.UserUrl+"/") {
		http.NotFound(w, r)
		return
	}

	proxyUrl := strings.TrimPrefix(path, setup.ConfigData.Server.UserUrl+"/")

	urlData, urlModel := users.GetUrlData(proxyUrl)

	if urlData == "" {
		http.NotFound(w, r)
		return
	}

	if urlModel == "html" {
		w.Header().Set("Content-Type", "text/html")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	//w.Write([]byte(urlData))
	fmt.Fprint(w, urlData)
}

func (s Server) download(w http.ResponseWriter, r *http.Request) {

	fileName := strings.TrimPrefix(r.URL.Path, setup.ConfigData.Download.Url+"/")
	filePath := filepath.Join(setup.ConfigData.Download.Folder, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	http.ServeFile(w, r, filePath)
}

func getConfigData() error {

	if setup.ConfigData.Proxy.OnlyName != "" {
		proxy.OnlyName = setup.ConfigData.Proxy.OnlyName
	} else {
		proxy.OnlyName, _ = random.GenerateStrings(6)
	}

	if setup.ConfigData.Download.Enabled {
		isRun := true

		if read.CreateFolder(setup.ConfigData.Download.Folder) != nil {
			isRun = false
		}

		if isRun {
			if read.CreateFolder(setup.ConfigData.Download.Folder+"/app") != nil {
				isRun = false
			}
		}

		if isRun {
			if read.CreateFolder(setup.ConfigData.Download.Folder+"/rules") != nil {
				isRun = false
			}
		}

		if isRun {
			download.DownloadRules()
			download.DownloadApp()
			download.RulesScheduledTasks()
			go download.AppScheduledTasks()
		}
	}

	protocol.GetSBString()
	if setup.ConfigData.Backup.Enabled {
		backup.GetProxyUrl()
		go backup.GetUrlTicker()
	}

	//fmt.Println(backup.ProxySBData)
	//fmt.Println(backup.ProxyUrlData)
	//fmt.Println(backup.SBSelectorOrUrlTestData)

	if setup.ConfigData.Proxy.Core == "sing-box" {
		proxy.ConfigData = singbox.Config{}
		proxy.LConfigData = singbox.LConfig{}
	} else {
		proxy.ConfigData = xray.Config{}
		proxy.LConfigData = singbox.LConfig{}
	}

	if read.CheckExistence(setup.ConfigData.Users.Config) != "file" {
		proxy.ConfigData.RenewData("new")

	} else {
		var config users.Config
		if err := read.GetJsonData(setup.ConfigData.Users.Config, &config); err == nil {
			users.ConfigData = config
		}
	}

	return nil

}

func (s Server) Run() {
	err := getConfigData()
	if err != nil {
		return
	}

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.Handle(setup.ConfigData.Server.Home.Url+"/sub-ui", http.HandlerFunc(s.home))

	http.Handle(setup.ConfigData.Server.Home.Url+"/login", http.HandlerFunc(s.login))

	http.Handle(setup.ConfigData.Server.Home.Url+"/logout", http.HandlerFunc(s.logout))

	http.Handle(setup.ConfigData.Server.UserUrl+"/", http.HandlerFunc(s.usersUrl))

	http.Handle(setup.ConfigData.Server.Post.Set, http.HandlerFunc(s.setTag))

	http.Handle(setup.ConfigData.Server.Post.Renew, http.HandlerFunc(s.renewUrl))

	if setup.ConfigData.Download.Enabled {
		http.Handle(setup.ConfigData.Download.Url+"/", http.HandlerFunc(s.download))
	}

	if setup.ConfigData.Backup.Enabled {
		http.Handle(setup.ConfigData.Server.Post.Backup, http.HandlerFunc(s.renewBackupUrl))
	}

	fmt.Println("启动服务器 :", setup.ConfigData.Server.Port)
	if err := http.ListenAndServe(":"+setup.ConfigData.Server.Port, nil); err != nil {
		fmt.Println("服务无法启动:", err)
	}

}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{
		"status":  "error",
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}
