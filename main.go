package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type CarResponse struct {
	Name string  `json:"name"`
	Dist float64 `json:"dist_km"`
}

type AddCarRequest struct {
	Plate string  `json:"plate" binding:"required"`
	Lng   float64 `json:"lng" binding:"required"`
	Lat   float64 `json:"lat" binding:"required"`
}

func main() {
	// ---------- 1. 连接MySQL ----------
	dsn := "root:Root@123456@tcp(127.0.0.1:3306)/car_db?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("MySQL连接失败: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		panic("MySQL无法ping通: " + err.Error())
	}
	fmt.Println("✅ MySQL连接成功")

	// ---------- 2. 连接Redis ----------
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
	fmt.Println("✅ Redis连接成功")

	// ---------- 3. 从MySQL加载车辆到Redis ----------
	rdb.Del(ctx, "cars")
	rows, err := db.QueryContext(ctx, "SELECT plate, lng, lat FROM cars WHERE status = 1")
	if err != nil {
		panic("查询车辆失败: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var plate string
		var lng, lat float64
		rows.Scan(&plate, &lng, &lat)
		rdb.GeoAdd(ctx, "cars", &redis.GeoLocation{
			Name:      plate,
			Longitude: lng,
			Latitude:  lat,
		})
		count++
	}
	fmt.Printf("✅ 从MySQL加载了 %d 辆空闲车到Redis\n", count)

	// ---------- 4. 创建Web服务 ----------
	r := gin.Default()

	// 接口1：查询附近车辆（GET）
	r.GET("/nearby", func(c *gin.Context) {
		lngStr := c.Query("lng")
		latStr := c.Query("lat")
		if lngStr == "" || latStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请传入经纬度参数"})
			return
		}

		lng, err := strconv.ParseFloat(lngStr, 64)
		lat, err2 := strconv.ParseFloat(latStr, 64)
		if err != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "经纬度格式错误"})
			return
		}

		results, err := rdb.GeoRadius(ctx, "cars", lng, lat, &redis.GeoRadiusQuery{
			Radius:     20,
			Unit:       "km",
			WithDist:   true,
			Sort:       "ASC",
		}).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}

		var resp []CarResponse
		for _, res := range results {
			resp = append(resp, CarResponse{
				Name: res.Name,
				Dist: res.Dist,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    resp,
		})
	})

	// 接口2：添加车辆（POST）
	r.POST("/car", func(c *gin.Context) {
		var req AddCarRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误，请提供plate、lng、lat"})
			return
		}

		_, err := db.ExecContext(ctx,
			"INSERT INTO cars (plate, lng, lat, status) VALUES (?, ?, ?, 1)",
			req.Plate, req.Lng, req.Lat,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败，可能车牌已存在: " + err.Error()})
			return
		}

		rdb.GeoAdd(ctx, "cars", &redis.GeoLocation{
			Name:      req.Plate,
			Longitude: req.Lng,
			Latitude:  req.Lat,
		})

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("✅ 车辆 %s 入库并同步缓存成功", req.Plate),
			"data":    req,
		})
	})

	// ---------- 接口3：删除车辆（DELETE）新增 ----------
	r.DELETE("/car/:plate", func(c *gin.Context) {
		plate := c.Param("plate")
		if plate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请提供车牌号"})
			return
		}

		result, err := db.ExecContext(ctx, "DELETE FROM cars WHERE plate = ?", plate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到该车牌: " + plate})
			return
		}

		// 从Redis中移除
		err = rdb.ZRem(ctx, "cars", plate).Err()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": fmt.Sprintf("⚠️ MySQL已删除，但Redis缓存清理失败: %s", plate),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("✅ 车辆 %s 已从数据库和缓存中删除", plate),
		})
	})

	// ---------- 5. 启动服务 ----------
	fmt.Println("🚀 服务已启动，数据持久化已开启（MySQL + Redis）")
	fmt.Println("📍 查询: GET http://localhost:8080/nearby?lng=116.42&lat=39.90")
	fmt.Println("➕ 添加: POST http://localhost:8080/car (JSON)")
	fmt.Println("➖ 删除: DELETE http://localhost:8080/car/{车牌号}")
	r.Run(":8080")
}