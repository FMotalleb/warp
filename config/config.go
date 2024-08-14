package config

import "time"

type Config struct {
	ListenProto     string
	ListenAddr      string
	ListenPort      uint16
	RemoteProto     string
	RemoteAddr      string
	RemotePort      uint16
	Threads         uint16
	Timeout         time.Duration
	Intercept       bool
	Base64Intercept bool
}
