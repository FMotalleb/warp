package config

import "time"

type GlobalConfig struct {
	Threads         uint16
	Timeout         time.Duration
	Intercept       bool
	Base64Intercept bool
}
