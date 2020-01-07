package main

import (
	"github.com/kataras/iris/v12"
	"tinyurl/controller"
)

func main() {
	app := iris.Default()
	app.Get("/{url:string regexp(^[0-9A-Za-z\\-_]{6})}", controller.GetUrl)
	app.Post("/add_url", controller.AddUrl)
	app.Run(iris.Addr(":4396"))
}