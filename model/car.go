package model

import "time"

type Car struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Plate     string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"plate"`
    Lng       float64   `gorm:"type:decimal(10,7);not null" json:"lng"`
    Lat       float64   `gorm:"type:decimal(10,7);not null" json:"lat"`
    Status    int8      `gorm:"default:1" json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

func (Car) TableName() string {
    return "cars"
}