package xray

type Client struct {
	Email    string `json:"email"`
	Id       string `json:"id"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Flow     string `json:"flow"`
}

type Fallback struct {
	Dest interface{} `josn:"dest"`
}

type Settings struct {
	Method    string     `json:"method"`
	Password  string     `json:"password"`
	Clients   []Client   `jsong:"clients"`
	Fallbacks []Fallback `json:"fallbacks"`
}

type XhttpSettings struct {
	Path string `json:"path"`
}

type HttpSettings struct {
	Host []string `json:"host"`
	Path string   `json:"path"`
}

type GrpcSettings struct {
	ServiceName string `json:"serviceName"`
}

type WsSettings struct {
	Path string `json:"path"`
}

type HttpupgradeSettings struct {
	Path string `json:"path"`
}

type SplithttpSettings struct {
	Path string `json:"path"`
}

type RealitySettings struct {
	ServerNames []string `json:"serverNames"`
	PrivateKey  string   `json:"privateKey"`
	ShortIds    []string `json:"shortIds"`
}

type TlsSettings struct {
	Alpn []string `json:"alpn"`
}

type StreamSettings struct {
	Network             string              `json:"network"`
	Security            string              `json:"security"`
	HttpSettings        HttpSettings        `json:"httpSettings"`
	GrpcSettings        GrpcSettings        `json:"grpcSettings"`
	WsSettings          WsSettings          `json:"wsSettings"`
	HttpupgradeSettings HttpupgradeSettings `json:"httpupgradeSettings"`
	SplithttpSettings   SplithttpSettings   `json:"splithttpSettings"`
	XhttpSettings       XhttpSettings       `json:"xhttpSettings"`
	RealitySettings     RealitySettings     `json:"realitySettings"`
	TlsSettings         TlsSettings         `json:"tlsSettings"`
}

type Inbound struct {
	Tag            string         `json:"tag"`
	Listen         string         `json:"listen"`
	Port           uint16         `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
}

type Config struct {
	Inbounds []Inbound `json:"inbounds"`
}

type RealityFallback struct {
	Index int
	Dest  interface{}
}

type RealityFallbacks struct {
	Fallbacks []RealityFallback
}
