package structs

type Peer struct {
	Name        string
	PublicKey   string
	PrivateKey  string
	Address     string
	Endpoint    string
	ConnectedTo []string
	Lighthouse  bool
	AllowedIps  string
}

type WgConfig struct {
	Peers               []Peer
	PersistentKeepAlive uint16
	Dns                 string
}
