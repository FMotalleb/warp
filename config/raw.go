package config

type RawConfig struct {
	GlobalConfig
	ListenProto string
	ListenAddr  string
	ListenPort  uint16
	RemoteProto string
	RemoteAddr  string
	RemotePort  uint16
}
