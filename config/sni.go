package config

type SNIConfig struct {
	GlobalConfig
	ListenAddr string
	ListenPort uint16
	RemotePort uint16
}
