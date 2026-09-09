package config

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "nearby_cars/model"
)

var DB *gorm.DB

func InitDB(dsn string) error {
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }
    // 自动迁移
    DB.AutoMigrate(&model.Car{})
    return nil
}