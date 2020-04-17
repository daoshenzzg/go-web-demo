package breaker

import (
	"github.com/afex/hystrix-go/hystrix"
	xtime "go-web-demo/library/time"
	"sync"
	"time"
)

// Config broker config.
type Config struct {
	// The namespace of the service with optional protocol prefix.
	// A namespace must be specified for all clients.
	Namespace              string
	Timeout                xtime.Duration
	MaxConcurrentRequests  int
	RequestVolumeThreshold uint64
	SleepWindow            xtime.Duration
	ErrorPercentThreshold  int
}

// Group represents a class of CircuitBreaker and forms a namespace in which
// units of CircuitBreaker.
type Group struct {
	mu        sync.RWMutex
	namespace string
	settings  map[string]bool
	conf      *Config
}

var (
	_mu   sync.RWMutex
	_conf = &Config{
		Namespace:              "default",
		Timeout:                xtime.Duration(hystrix.DefaultTimeout),
		MaxConcurrentRequests:  hystrix.DefaultMaxConcurrent,
		RequestVolumeThreshold: uint64(hystrix.DefaultVolumeThreshold),
		SleepWindow:            xtime.Duration(hystrix.DefaultSleepWindow),
		ErrorPercentThreshold:  hystrix.DefaultErrorPercentThreshold,
	}
)

func (conf *Config) fix() {
	if conf.Namespace == "" {
		conf.Namespace = "default"
	}
	if conf.Timeout <= 0 {
		conf.Timeout = xtime.Duration(hystrix.DefaultTimeout)
	}
	if conf.MaxConcurrentRequests <= 0 {
		conf.MaxConcurrentRequests = hystrix.DefaultMaxConcurrent
	}
	if conf.RequestVolumeThreshold <= 0 {
		conf.RequestVolumeThreshold = uint64(hystrix.DefaultVolumeThreshold)
	}
	if conf.SleepWindow == 0 {
		conf.SleepWindow = xtime.Duration(hystrix.DefaultSleepWindow)
	}
	if conf.ErrorPercentThreshold <= 0 {
		conf.ErrorPercentThreshold = hystrix.DefaultErrorPercentThreshold
	}
}

// NewGroup new a breaker group container, if conf nil use default conf.
func NewGroup(conf *Config) *Group {
	if conf == nil {
		_mu.RLock()
		conf = _conf
		_mu.RUnlock()
	} else {
		conf.fix()
	}
	return &Group{
		namespace: conf.Namespace,
		settings:  make(map[string]bool),
		conf:      conf,
	}
}

// Reload reload the group by specified config, this may let all inner breaker
// reset to a new one.
func (g *Group) Reload(conf *Config) {
	if conf == nil {
		return
	}
	conf.fix()
	g.mu.Lock()
	g.conf = conf
	g.mu.Unlock()
}

// Warped name with namespace for Hystrix DO.
func (g *Group) Do(name string, run func() error) (err error) {
	name = g.namespace + "-" + name
	g.setBreakerConfig(name)
	return hystrix.Do(name, func() error {
		return run()
	}, nil)
}

// setBreakerConfig set breaker configuration atomic if not set
func (g *Group) setBreakerConfig(name string) {
	if _, ok := g.settings[name]; !ok {
		g.mu.Lock()
		defer g.mu.Unlock()

		if _, ok := g.settings[name]; !ok {
			hystrix.ConfigureCommand(name, hystrix.CommandConfig{
				Timeout:                int(time.Duration(g.conf.Timeout) / time.Millisecond),
				MaxConcurrentRequests:  g.conf.MaxConcurrentRequests,
				RequestVolumeThreshold: int(g.conf.RequestVolumeThreshold),
				SleepWindow:            int(time.Duration(g.conf.SleepWindow) / time.Millisecond),
				ErrorPercentThreshold:  g.conf.ErrorPercentThreshold,
			})

			copy := make(map[string]bool)
			for key, val := range g.settings {
				copy[key] = val
			}
			copy[name] = true
			g.settings = copy
		}
	}
}
