package model

import (
	"time"
)

const TableNamePoint = "point"

// Point mapped from table <point>
type Point struct {
	PoID         uint      `gorm:"type:INTEGER;primaryKey;autoIncrement"`
	MbID         string    `gorm:"type:VARCHAR(20);not null;default:''"`
	Member       Member    `gorm:"foreignKey:MbID;not null;default:''"`
	PoDatetime   time.Time `gorm:"type:DATETIME;not null;autoCreateTime"`
	PoContent    string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoint      int       `gorm:"type:INTEGER;not null;default:0"`
	PoUsePoint   int       `gorm:"type:INTEGER;not null;default:0"`
	PoExpired    int       `gorm:"type:INTEGER;not null;default:0"`
	PoExpireDate time.Time `gorm:"type:DATE;not null;autoCreateTime"`
	PoMbPoint    int       `gorm:"type:INTEGER;not null;default:0"`
	PoRelTable   string    `gorm:"type:VARCHAR(20);not null;default:''"`
	PoRelID      string    `gorm:"type:VARCHAR(20);not null;default:''"`
	PoRelAction  string    `gorm:"type:VARCHAR(100);not null;default:''"`
}

// TableName Point's table name
func (*Point) TableName() string {
	return Prefix + TableNamePoint
}
