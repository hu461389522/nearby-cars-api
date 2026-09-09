package service

import (
    "context"
    "fmt"
    "nearby_cars/config"
    "nearby_cars/model"

    "github.com/go-redis/redis/v8"
   
)

// 查询附近车辆（从Redis读取）
func GetNearbyCars(ctx context.Context, rdb *redis.Client, lng, lat float64) ([]map[string]interface{}, error) {
    results, err := rdb.GeoRadius(ctx, "cars", lng, lat, &redis.GeoRadiusQuery{
        Radius:    20,
        Unit:      "km",
        WithDist:  true,
        Sort:      "ASC",
    }).Result()
    if err != nil {
        return nil, err
    }

    var resp []map[string]interface{}
    for _, res := range results {
        resp = append(resp, map[string]interface{}{
            "name":    res.Name,
            "dist_km": res.Dist,
        })
    }
    return resp, nil
}

// 添加车辆（同时写入MySQL和Redis）
func AddCar(ctx context.Context, rdb *redis.Client, plate string, lng, lat float64) error {
    car := model.Car{
        Plate:  plate,
        Lng:    lng,
        Lat:    lat,
        Status: 1,
    }
    // 写入MySQL
    if err := config.DB.Create(&car).Error; err != nil {
        return fmt.Errorf("MySQL写入失败: %w", err)
    }
    // 写入Redis
    rdb.GeoAdd(ctx, "cars", &redis.GeoLocation{
        Name:      plate,
        Longitude: lng,
        Latitude:  lat,
    })
    return nil
}

// 删除车辆
func DeleteCar(ctx context.Context, rdb *redis.Client, plate string) error {
    // 从MySQL删除
    result := config.DB.Where("plate = ?", plate).Delete(&model.Car{})
    if result.RowsAffected == 0 {
        return fmt.Errorf("未找到车牌: %s", plate)
    }
    if result.Error != nil {
        return result.Error
    }
    // 从Redis删除
    rdb.ZRem(ctx, "cars", plate)
    return nil
}

// 加载所有车辆到Redis（服务启动时调用）
func LoadCarsToRedis(ctx context.Context, rdb *redis.Client) (int, error) {
    var cars []model.Car
    if err := config.DB.Where("status = ?", 1).Find(&cars).Error; err != nil {
        return 0, err
    }
    // 清空Redis
    rdb.Del(ctx, "cars")
    for _, car := range cars {
        rdb.GeoAdd(ctx, "cars", &redis.GeoLocation{
            Name:      car.Plate,
            Longitude: car.Lng,
            Latitude:  car.Lat,
        })
    }
    return len(cars), nil
}