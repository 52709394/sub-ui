package users

type User struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
	Flow     string `json:"flow"`
	UserPath string `json:"user_path"`
}

type Inbound struct {
	Addr              string `json:"address"`
	Protocol          string `json:"protocol"`
	Port              uint16 `json:"port"`
	Users             []User `json:"users"`
	CongestionControl string `json:"congestion_control"`
	Network           string `json:"network"`
	Host              string `json:"host"`
	Path              string `json:"path"`
	ServiceName       string `json:"service_name"`
	Security          string `json:"security"`
	FixedSecurity     bool   `json:"fixed_Security"`
	Alpn              string `json:"alpn"`
	Sni               string `json:"sni"`
	PublicKey         string `json:"publicKey"`
	ShortId           string `json:"shortId"`
	Fingerprint       string `json:"fingerprint"`
	Tag               string `json:"tag"`
	TagPath           string `json:"tag_path"`
}

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}

var ConfigData Config

type ResetUrl struct {
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
