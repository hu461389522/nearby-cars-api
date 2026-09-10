package router

import (
	"github.com/gin-gonic/gin"
	"nearby_cars/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 静态文件
	r.Static("/static", "./static")

	// API 接口
	r.GET("/nearby", controller.GetNearbyCars)
	r.POST("/car", controller.AddCar)
	r.DELETE("/car/:plate", controller.DeleteCar)
	r.POST("/chat", controller.Chat)

	// 首页跳转到聊天页面
	r.GET("/", func(c *gin.Context) {
		c.File("./static/chat.html")
	})

	return r
}
