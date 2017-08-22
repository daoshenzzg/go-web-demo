package mvc

import (
	"net/http"
	"runtime"
	"time"

	"person.mgtv.com/framework/config"
	"person.mgtv.com/framework/logs"
)

type App struct {
	Handlers *ControllerHandler
	Server   *http.Server
}

func NewApp() *App {
	c := NewControllerHandler()
	app := &App{Handlers: c, Server: &http.Server{}}
	return app
}

func (app *App) Router(name string, controller ControllerInterface) {
	app.Handlers.Add(name, controller)
}

func (app *App) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Server.Addr = ":" + config.ServerPort
	app.Server.Handler = app.Handlers
	app.Server.ReadTimeout = time.Duration(config.ServerTimeout) * time.Millisecond
	app.Server.WriteTimeout = time.Duration(config.ServerTimeout) * time.Millisecond

	endRunning := make(chan bool, 1)

	logs.GetLogger("system").Infof("Http server running on %s", config.ServerPort)

	// 监听端口
	if err := app.Server.ListenAndServe(); err != nil {
		logs.GetLogger("system").Errorf("Http server start error: %s", err.Error())
	}

	<-endRunning
}
