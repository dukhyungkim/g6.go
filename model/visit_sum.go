package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameVisitSum = "visit_sum"

type VisitSum struct {
	VsDate  string `gorm:"type:DATE;primaryKey;not null;default:''" json:"vs_date"`
	VsCount int64  `gorm:"type:INTEGER;not null;default:0" json:"vs_count"`
}

func (*VisitSum) TableName() string {
	return config.Global.DbTablePrefix + TableNameVisitSum
}
