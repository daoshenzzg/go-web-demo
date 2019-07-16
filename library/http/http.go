package http

import (
	xtime "go-web-demo/library/time"
	"net"
	"net/http"
	"time"
)

// Config client settings.
type Config struct {
	MaxTotal    int
	MaxPerHost  int
	KeepAlive   xtime.Duration
	DialTimeout xtime.Duration
	Timeout     xtime.Duration
}

// NewClient new http client pool
func NewClient(c *Config) *http.Client {
	if c.DialTimeout <= 0 || c.Timeout <= 0 {
		panic("must config http timeout")
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(c.DialTimeout),
			}).DialContext,
			MaxIdleConns:        c.MaxTotal,
			MaxIdleConnsPerHost: c.MaxPerHost,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: time.Duration(c.Timeout),
	}
}
