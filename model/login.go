package model

import (
	"github.com/dukhyungkim/gonuboard/config"
	"time"
)

const TableNameLogin = "login"

type Login struct {
	LoID       uint      `gorm:"type:INTEGER;primaryKey;autoincrement"`
	LoIP       string    `gorm:"type:VARCHAR(100);not null;default:''"`
	MbID       string    `gorm:"type:VARCHAR(20);not null;default:''"`
	LoDatetime time.Time `gorm:"type:DATETIME;not null;default:''"`
	LoLocation string    `gorm:"type:TEXT;not null;default:''"`
	LoURL      string    `gorm:"type:TEXT;not null;default:''"`
}

func (*Login) TableName() string {
	return config.Global.DbTablePrefix + TableNameLogin
}
