package structs

type Peer struct {
	name       string
	publicKey  string
	privateKey string
	address    string
	endpoint   string
}

type WgConfig struct {
	peer                []Peer
	persistentKeepAlive uint16
}
