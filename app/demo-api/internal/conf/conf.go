package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"go-web-demo/library/cache/redis"
	"go-web-demo/library/database/sql"
	"go-web-demo/library/log"
	xhttp "go-web-demo/library/net/http"
	xtime "go-web-demo/library/time"
)

var (
	httpPort int
	confPath string
	Conf     = &Config{}
)

type Config struct {
	// App
	App *App
	// Log
	Log *log.Config
	// DB
	MySQL *MySQL
	// Redis
	Redis *redis.Config
	// HttpClient
	HttpClient map[string]*HttpClient
}

type App struct {
	HttpPort     int
	RunMode      string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

type MySQL struct {
	School *sql.Config
}

type HttpClient struct {
	Addr       string
	ClientConf *xhttp.ClientConfig
}

func init() {
	flag.IntVar(&httpPort, "http.port", -1, "http port")
	flag.StringVar(&confPath, "conf", "./app/demo-api/configs/application.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
