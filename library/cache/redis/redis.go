package redis

import (
	"github.com/gomodule/redigo/redis"
	xtime "go-web-demo/library/time"
	"time"
)

// Pool.
type Pool struct {
	redis.Pool
}

// Config client settings.
type Config struct {
	Addr         string
	MaxIdle      int
	MaxActive    int
	IdleTimeout  xtime.Duration
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// NewPool creates a new pool.
func NewPool(c *Config) (p *Pool) {
	if c.DialTimeout <= 0 || c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
		panic("must config redis timeout")
	}

	dialFunc := func() (redis.Conn, error) {
		return redis.Dial(
			"tcp",
			c.Addr,
			redis.DialConnectTimeout(time.Duration(c.DialTimeout)),
			redis.DialReadTimeout(time.Duration(c.ReadTimeout)),
			redis.DialWriteTimeout(time.Duration(c.WriteTimeout)))
	}

	return &Pool{redis.Pool{
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial:        dialFunc}}
}

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }
