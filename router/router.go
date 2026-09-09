package router

import (
    "github.com/gin-gonic/gin"
    "nearby_cars/controller"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    r.GET("/nearby", controller.GetNearbyCars)
    r.POST("/car", controller.AddCar)
    r.DELETE("/car/:plate", controller.DeleteCar)
    
    return r
}