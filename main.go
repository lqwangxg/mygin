package main

import (
	"mygin/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/static", "./static")

	//-------------------------------------------
	engine.GET("/", controller.Index)
	engine.GET("/match", controller.Match)
	engine.POST("/match", controller.Match)
	engine.GET("/show", controller.Show)
	engine.GET("/sysdate", controller.SysDate)
	engine.GET("/query", controller.Exec)
	engine.POST("/query", controller.Exec)
	//-------------------------------------------

	engine.Run(":3000")
}
