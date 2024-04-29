package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameContent = "content"

type Content struct {
	CoID            string `gorm:"type:VARCHAR(20);not null;default:''"`
	CoHTML          int    `gorm:"type:INTEGER;not null;default:0"`
	CoSubject       string `gorm:"type:VARCHAR(255);not null;default:''"`
	CoContent       string `gorm:"type:TEXT;not null;default:''"`
	CoSeoTitle      string `gorm:"type:VARCHAR(255);not null;default:''"`
	CoMobileContent string `gorm:"type:TEXT;not null;default:''"`
	CoSkin          string `gorm:"type:VARCHAR(255);not null;default:''"`
	CoMobileSkin    string `gorm:"type:VARCHAR(255);not null;default:''"`
	CoTagFilterUse  int    `gorm:"type:INTEGER;not null;default:0"`
	CoHit           int    `gorm:"type:INTEGER;not null;default:0"`
	CoIncludeHead   string `gorm:"type:VARCHAR(255);not null;default:''"`
	CoIncludeTail   string `gorm:"type:VARCHAR(255);not null;default:''"`
}

func (*Content) TableName() string {
	return config.Global.DbTablePrefix + TableNameContent
}
