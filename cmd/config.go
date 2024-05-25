package cmd

import "net/url"

type Config struct {
	Listen  *url.URL
	Remote  *url.URL
	Threads uint16
}
