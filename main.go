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
	engine.GET("/exec", controller.Exec)
	engine.POST("/exec", controller.Exec)
	engine.GET("/getvalue", controller.QueryValue)
	engine.POST("/getvalue", controller.QueryValue)
	engine.GET("/gettable", controller.QueryTable)
	engine.POST("/gettable", controller.QueryTable)

	//-------------------------------------------

	engine.Run(":3000")
}
