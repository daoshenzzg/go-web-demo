package main

import (
	demoController "person.mgtv.com/controller/demo"
	"person.mgtv.com/framework/mvc"
)

func main() {
	app := mvc.NewApp()

	app.Router("demo", &demoController.DemoController{})

	app.Run()
}
