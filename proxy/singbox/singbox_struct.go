package singbox

type User struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
	Flow     string `json:"flow"`
}

type Handshake struct {
	Server string `json:"server"`
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
	Tag               string    `json:"tag"`
	Type              string    `json:"type"`
	Listen            string    `json:"listen"`
	Port              uint16    `json:"listen_port"`
	Method            string    `json:"method"`
	Password          string    `json:"password"`
	Users             []User    `json:"users"`
	Version           uint16    `json:"version"`
	Detour            string    `json:"detour"`
	Handshake         Handshake `json:"handshake"`
	CongestionControl string    `json:"congestion_control"`
	Transport         Transport `json:"transport"`
	Tls               Tls       `json:"tls"`
}

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}

type LUser struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

type LInbound struct {
	Tag      string  `json:"tag"`
	Type     string  `json:"type"`
	Method   string  `json:"method"`
	Password string  `json:"password"`
	Users    []LUser `json:"users"`
}

type LConfig struct {
	Inbounds []LInbound `json:"inbounds"`
}

type Detour struct {
	Index  int
	Detour string
}

type SBDetours struct {
	Detours []Detour
}
