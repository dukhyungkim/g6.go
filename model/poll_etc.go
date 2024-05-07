package model

import (
	"github.com/dukhyungkim/gonuboard/config"
	"time"
)

const TableNamePollEtc = "poll_etc"

type PollEtc struct {
	PcID       uint `gorm:"type:INTEGER;primaryKey;autoIncrement"`
	PoID       uint
	Poll       Poll      `gorm:"type:foreignKey;PoID;not null;default:0"`
	MbID       string    `gorm:"type:VARCHAR(20);not null;default:''"`
	PcName     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PcIdea     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	PcDatetime time.Time `gorm:"type:DATETIME;not null;autoCreateTime"`
}

func (*PollEtc) TableName() string {
	return config.Global.DbTablePrefix + TableNamePollEtc
}
