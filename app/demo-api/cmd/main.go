package main

import (
	"context"
	"flag"
	"go-web-demo/app/demo-api/internal/conf"
	"go-web-demo/app/demo-api/internal/server/http"
	"go-web-demo/app/demo-api/internal/service"
	"go-web-demo/library/log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	flag.Parse()

	// IDE中，你也可以放开注释，直接运行。
	dir, _ := filepath.Abs("./app/demo-api/configs/application.toml")
	flag.Set("conf", dir)

	// 初始化配置
	conf.Init()
	// 初始化日志
	log.Init(conf.Conf.Log)
	defer log.Close()
	srv := service.New(conf.Conf)
	log.Info("go-web-demo start")
	httpSrv := http.New(conf.Conf, srv)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Error("httpSrv.Shutdown error(%v)", err)
			}
			log.Info("go-web-demo exit")
			httpSrv.Close()
			cancel()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
