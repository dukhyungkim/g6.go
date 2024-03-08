package model

import (
	"time"
)

const TableNameLogin = "login"

// Login mapped from table <login>
type Login struct {
	LoID       uint      `gorm:"type:INTEGER;primaryKey;autoincrement"`
	LoIP       string    `gorm:"type:VARCHAR(100);not null;default:''"`
	MbID       string    `gorm:"type:VARCHAR(20);not null;default:''"`
	LoDatetime time.Time `gorm:"type:DATETIME;not null;default:''"`
	LoLocation string    `gorm:"type:TEXT;not null;default:''"`
	LoURL      string    `gorm:"type:TEXT;not null;default:''"`
}

// TableName Login's table name
func (*Login) TableName() string {
	return Prefix + TableNameLogin
}
