package main

import (
	"context"
	"fmt"
	"nearby_cars/config"
	"nearby_cars/controller"
	"nearby_cars/router"
	"nearby_cars/service"

	"github.com/go-redis/redis/v8"
)

func main() {
	// 1. 初始化MySQL
	dsn := "root:123456@tcp(127.0.0.1:3306)/car_db?charset=utf8mb4&parseTime=True"
	if err := config.InitDB(dsn); err != nil {
		panic("MySQL连接失败: " + err.Error())
	}
	fmt.Println("✅ MySQL连接成功")

	// 2. 初始化Redis
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

	// 3. 设置controller中的Redis客户端
	controller.InitRedis(rdb)

	// 4. 加载数据到Redis
	count, err := service.LoadCarsToRedis(ctx, rdb)
	if err != nil {
		panic("加载数据到Redis失败: " + err.Error())
	}
	fmt.Printf("✅ 加载了 %d 辆车到Redis\n", count)

	// 5. 启动服务
	r := router.SetupRouter()
	fmt.Println("🚀 服务启动（分层架构）")
	r.Run(":8080")
}