package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameVisit = "visit"

type Visit struct {
	ViID      uint   `gorm:"type:INTEGER;primaryKey;autoIncrement"`
	ViIP      string `gorm:"type:VARCHAR(100);not null;default:''"`
	ViDate    string `gorm:"type:DATE;not null;default:''"`
	ViTime    string `gorm:"type:TIME;not null;default:''"`
	ViReferer string `gorm:"type:TEXT;not null;default:''"`
	ViAgent   string `gorm:"type:VARCHAR(200);not null;default:''"`
	ViBrowser string `gorm:"type:VARCHAR(255);not null;default:''"`
	ViOs      string `gorm:"type:VARCHAR(255);not null;default:''"`
	ViDevice  string `gorm:"type:VARCHAR(255);not null;default:''"`
}

func (*Visit) TableName() string {
	return config.Global.DbTablePrefix + TableNameVisit
}
