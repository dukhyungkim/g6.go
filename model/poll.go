package model

import (
	"github.com/dukhyungkim/gonuboard/config"
	"time"
)

const TableNamePoll = "poll"

type Poll struct {
	PoID      uint      `gorm:"type:INTEGER;primaryKey;autoIncrement"`
	PoSubject string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll1   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll2   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll3   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll4   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll5   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll6   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll7   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll8   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoPoll9   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoCnt1    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt2    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt3    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt4    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt5    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt6    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt7    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt8    int       `gorm:"type:INTEGER;not null;default:0"`
	PoCnt9    int       `gorm:"type:INTEGER;not null;default:0"`
	PoEtc     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PoLevel   int       `gorm:"type:INTEGER;not null;default:0"`
	PoPoint   int       `gorm:"type:INTEGER;not null;default:0"`
	PoDate    time.Time `gorm:"type:DATE;not null;autoCreateTime"`
	PoIps     string    `gorm:"type:TEXT;not null;default:''"`
	MbIds     string    `gorm:"type:TEXT;not null;default:''"`
	PoUse     int       `gorm:"type:INTEGER;not null;default:0"`
}

func (*Poll) TableName() string {
	return config.Global.DbTablePrefix + TableNamePoll
}
