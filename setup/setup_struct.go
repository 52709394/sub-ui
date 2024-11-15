package setup

var CookieName string
var CookieValue string
var CookieDay int

type ConstUser struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Consts struct {
	Tag   string      `json:"tag"`
	Users []ConstUser `json:"users"`
}

type Static struct {
	Enabled   bool     `json:"enabled"`
	ConstList []Consts `json:"const_list"`
}

type BckProxy struct {
	SBTag string `json:"sb_tag"`
	Url   string `json:"url"`
}

type BacSingBox struct {
	Outbound                 string `json:"outbound"`
	DownloadDetour           string `json:"download_detour"`
	ExternalUiDownloadDetour string `json:"external_ui_download_detour"`
}

type Exclude struct {
	Tag   string   `json:"tag"`
	Users []string `json:"users"`
}

type Backup struct {
	Enabled        bool       `json:"enabled"`
	StartTime      int        `json:"start_time"`
	UpdateInterval int        `json:"update_interval"`
	SBSelector     bool       `json:"sb_selector"`
	Excludes       []Exclude  `json:"excludes"`
	SingBox        BacSingBox `json:"sing-box"`
	ProxyList      []BckProxy `json:"proxy_list"`
}

type Github struct {
	User       string `json:"user"`
	Repository string `json:"repository"`
	Regexp     string `json:"regexp"`
	Name       string `json:"name"`
	Url        string `json:"url"`
}

type Rule struct {
	Url            string `json:"url"`
	Name           string `json:"name"`
	UpdateInterval int    `json:"update_interval"`
}

type Download struct {
	Enabled           bool     `json:"enabled"`
	StartTime         int      `json:"start_time"`
	AppUpdateInterval int      `json:"app_update_interval"`
	Folder            string   `json:"folder"`
	Url               string   `json:"url"`
	RuleList          []Rule   `json:"rule_list"`
	GithubList        []Github `json:"github_list"`
}

type SingBox struct {
	Config  string `json:"config"`
	MainTag string `json:"main_tag"`
}

type AppList struct {
	User       string `json:"user"`
	Repository string `json:"repository"`
	Regexp     string `json:"regexp"`
	OnlyCopy   bool   `json:"only_copy"`
	Label      string `json:"label"`
	Url        string `json:"url"`
}

type App struct {
	GitHubProxy string    `json:"github_proxy"`
	AppList     []AppList `json:"app_list"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Day   int    `json:"day"`
}

type Post struct {
	Set    string `json:"set"`
	Renew  string `json:"renew"`
	Backup string `json:"backup"`
}

type Home struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Server struct {
	Port    string `json:"port"`
	UserUrl string `json:"user_url"`
	Home    Home   `json:"home"`
	Post    Post   `json:"post"`
	Cookie  Cookie `json:"cookie"`
}

type Proxy struct {
	Core     string `json:"core"`
	Config   string `json:"config"`
	RealTime bool   `json:"real_time"`
	OnlyName string `json:"only_name"`
}

type Users struct {
	Domain     string `json:"domain"`
	Port       string `json:"port"`
	Config     string `json:"config"`
	UtlsFp     string `json:"utls_fp"`
	VmessModel string `json:"vmess_model"`
	Ws0Rtt     string `json:"ws_0-rtt"`
}

type Config struct {
	Users    Users    `json:"users"`
	Proxy    Proxy    `json:"proxy"`
	Server   Server   `json:"server"`
	SingBox  SingBox  `json:"sing-box"`
	Download Download `json:"download"`
	App      App      `json:"app"`
	Backup   Backup   `json:"backup"`
	Static   Static   `json:"static"`
}
