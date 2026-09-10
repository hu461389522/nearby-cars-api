package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"nearby_cars/service"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis(rdb *redis.Client) {
	Rdb = rdb
}

// GetNearbyCars 查询附近车辆
func GetNearbyCars(c *gin.Context) {
	lngStr := c.Query("lng")
	latStr := c.Query("lat")
	if lngStr == "" || latStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请传入经纬度参数"})
		return
	}
	lng, _ := strconv.ParseFloat(lngStr, 64)
	lat, _ := strconv.ParseFloat(latStr, 64)

	data, err := service.GetNearbyCars(Ctx, Rdb, lng, lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// AddCar 添加车辆
func AddCar(c *gin.Context) {
	var req struct {
		Plate string  `json:"plate" binding:"required"`
		Lng   float64 `json:"lng" binding:"required"`
		Lat   float64 `json:"lat" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}
	if err := service.AddCar(Ctx, Rdb, req.Plate, req.Lng, req.Lat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    req,
	})
}

// DeleteCar 删除车辆
func DeleteCar(c *gin.Context) {
	plate := c.Param("plate")
	if plate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供车牌号"})
		return
	}
	if err := service.DeleteCar(Ctx, Rdb, plate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// Chat 智能客服接口
func Chat(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供question参数"})
		return
	}

	answer, err := service.ChatWithAI(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    gin.H{"question": req.Question, "answer": answer},
	})
}