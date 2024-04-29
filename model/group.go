package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameGroup = "group"

type Group struct {
	GrID        string `gorm:"type:VARCHAR(10);primaryKey;not null;"`
	GrSubject   string `gorm:"type:VARCHAR(255);not null;default:''"`
	GrDevice    string `gorm:"type:VARCHAR(6);not null;default:both"`
	GrAdmin     string `gorm:"type:VARCHAR(255);not null;default:''"`
	GrUseAccess int    `gorm:"type:INTEGER;not null;default:0"`
	GrOrder     int    `gorm:"type:INTEGER;not null;default:0"`
	Gr1Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr2Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr3Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr4Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr5Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr6Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr7Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr8Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr9Subj     string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr10Subj    string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr1         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr2         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr3         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr4         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr5         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr6         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr7         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr8         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr9         string `gorm:"type:VARCHAR(255);not null;default:''"`
	Gr10        string `gorm:"type:VARCHAR(255);not null;default:''"`
}

func (*Group) TableName() string {
	return config.Global.DbTablePrefix + TableNameGroup
}
