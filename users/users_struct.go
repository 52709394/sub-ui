package users

type User struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid,omitempty"`
	Method   string `json:"method,omitempty"`
	Password string `json:"password,omitempty"`
	Flow     string `json:"flow,omitempty"`
	Static   bool   `json:"static"`
	UserPath string `json:"user_path"`
}

type Transport struct {
	Host        string `json:"host,omitempty"`
	Path        string `json:"path,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
}

type Tls struct {
	Sni  string `json:"sni,omitempty"`
	Alpn string `json:"alpn,omitempty"`
}

type Shadowtls struct {
	Version     uint16 `json:"version,omitempty"`
	Detour      string `json:"-"`
	Sni         string `json:"sni,omitempty"`
	DetourProxy string `json:"detour_proxy,omitempty"`
}

type Reality struct {
	Sni       string `json:"sni,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	ShortId   string `json:"shortId,omitempty"`
}

type Inbound struct {
	Addr              string     `json:"address"`
	Tag               string     `json:"tag"`
	TagPath           string     `json:"tag_path"`
	Protocol          string     `json:"protocol"`
	Port              uint16     `json:"port"`
	ServiceListen     string     `json:"-"`
	ServicePort       uint16     `json:"-"`
	Hide              bool       `json:"hide"`
	Users             []User     `json:"users"`
	CongestionControl string     `json:"congestion_control,omitempty"`
	Shadowtls         *Shadowtls `json:"shadowtls,omitempty"`
	Network           string     `json:"network"`
	Transport         *Transport `json:"transport,omitempty"`
	Security          string     `json:"security"`
	FixedSecurity     bool       `json:"fixed_Security"`
	Tls               *Tls       `json:"tls,omitempty"`
	Reality           *Reality   `json:"reality,omitempty"`
	Fingerprint       string     `json:"fingerprint"`
}

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}

var ConfigData Config

type RenewUsers struct {
	Mod   string     `json:"mod"`
	Users []UserData `json:"users"`
}

type BackupInfo struct {
	Mod   string     `json:"mod"`
	Users []UserData `json:"users"`
}

type UserData struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Name string `json:"name"`
}

type TagData struct {
	Tag      string `json:"tag"`
	Index    int    `json:"index"`
	Addr     string `json:"addr"`
	Port     uint16 `json:"port"`
	Security string `json:"security"`
	Alpn     string `json:"alpn"`
}
