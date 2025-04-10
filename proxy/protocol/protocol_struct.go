package protocol

type Config struct {
	Addr         *string
	Protocol     *string
	Port         *string
	TuicCC       *string
	Method       *string
	Network      *string
	Host         *string
	Path         *string
	ServiceName  *string
	Security     *string
	Alpn         *string
	Version      *string
	DetourProxy  *string
	Sni          *string
	PublicKey    *string
	ShortId      *string
	Fingerprint  *string
	UserName     *string
	UserUUID     *string
	UserPassword *string
	UserFlow     *string
}

type SBString struct {
	UrlTest          string
	Selector         string
	VmessWsTls       string
	VlessTcpReality  string
	VlessHttpReality string
	VlessGrpcReality string
	VlessTcpTls      string
	TrojanTcpTls     string
	Anytls           string
	Hysteria2        string
	Tuic             string
	Shadowtls        string
	ShadowtlsSS      string
	Shadowsocks      string
}

var SBStringData SBString
