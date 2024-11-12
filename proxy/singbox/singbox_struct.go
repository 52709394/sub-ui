package singbox

type User struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
	Flow     string `json:"flow"`
}

type Reality struct {
	Enabled    bool     `json:"enabled"`
	PrivateKey string   `json:"private_key"`
	ShortId    []string `json:"short_id"`
}

type Transport struct {
	Type        string   `json:"type"`
	Path        string   `json:"path"`
	Host        []string `json:"host"`
	ServiceName string   `json:"service_name"`
}

type Tls struct {
	Enabled    bool     `json:"enabled"`
	ServerName string   `json:"server_name"`
	Reality    Reality  `json:"reality"`
	Alpn       []string `json:"alpn"`
}

type Inbound struct {
	Type              string    `json:"type"`
	Listen            string    `json:"listen"`
	Port              uint16    `json:"listen_port"`
	Users             []User    `json:"users"`
	CongestionControl string    `json:"congestion_control"`
	Transport         Transport `json:"transport"`
	Tls               Tls       `json:"tls"`
	Tag               string    `json:"tag"`
}

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}
