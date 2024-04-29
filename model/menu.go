package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameMenu = "menu"

type Menu struct {
	MeID        int    `gorm:"type:INTEGER;primaryKey;autoincrement"`
	MeCode      string `gorm:"type:VARCHAR(255);not null;default:''"`
	MeName      string `gorm:"type:VARCHAR(255);not null;default:''"`
	MeLink      string `gorm:"type:VARCHAR(255);not null;default:''"`
	MeTarget    string `gorm:"type:VARCHAR(255);not null;default:''"`
	MeOrder     int    `gorm:"type:INTEGER;not null;default:0"`
	MeUse       int    `gorm:"type:INTEGER;not null;default:0"`
	MeMobileUse int    `gorm:"type:INTEGER;not null;default:0"`
}

func (*Menu) TableName() string {
	return config.Global.DbTablePrefix + TableNameMenu
}
