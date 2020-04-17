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
	"syscall"
	"time"
)

func main() {
	flag.Parse()

	// 初始化配置
	conf.Init()
	// 初始化日志
	log.Init(conf.Conf.Log)
	defer log.Close()
	srv := service.New(conf.Conf)
	httpSrv := http.New(conf.Conf, srv)
	log.Info("ab-test started, listening on port: %d, runMode: %s.",
		conf.Conf.App.HttpPort, conf.Conf.App.RunMode)
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
