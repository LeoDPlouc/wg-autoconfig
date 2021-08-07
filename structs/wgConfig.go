package structs

type Node struct {
	Name          string
	PublicKey     string
	PrivateKey    string
	Address       string
	Endpoint      string
	ConnectedTo   []string
	Lighthouse    bool
	AllowedIps    string
	PostUp        string
	PostDown      string
	ListeningPort string
}

type WgConfig struct {
	Nodes               []Node
	PersistentKeepAlive uint16
	Dns                 string
}
