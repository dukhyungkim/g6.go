package model

import "github.com/dukhyungkim/gonuboard/config"

const TableNameQaConfig = "qa_config"

type QaConfig struct {
	ID                  uint   `gorm:"type:INTEGER;primaryKey"`
	QaTitle             string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaCategory          string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaSkin              string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaMobileSkin        string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaUseEmail          int    `gorm:"type:INTEGER;not null;default:0"`
	QaReqEmail          int    `gorm:"type:INTEGER;not null;default:0"`
	QaUseHp             int    `gorm:"type:INTEGER;not null;default:0"`
	QaReqHp             int    `gorm:"type:INTEGER;not null;default:0"`
	QaUseSms            int    `gorm:"type:INTEGER;not null;default:0"`
	QaSendNumber        string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaAdminHp           string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaAdminEmail        string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaUseEditor         int    `gorm:"type:INTEGER;not null;default:0"`
	QaSubjectLen        int    `gorm:"type:INTEGER;not null;default:0"`
	QaMobileSubjectLen  int    `gorm:"type:INTEGER;not null;default:0"`
	QaPageRows          int    `gorm:"type:INTEGER;not null;default:0"`
	QaMobilePageRows    int    `gorm:"type:INTEGER;not null;default:0"`
	QaImageWidth        int    `gorm:"type:INTEGER;not null;default:0"`
	QaUploadSize        int    `gorm:"type:INTEGER;not null;default:0"`
	QaInsertContent     string `gorm:"type:TEXT;not null;default:''"`
	QaIncludeHead       string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaIncludeTail       string `gorm:"type:VARCHAR(255);not null;default:''"`
	QaContentHead       string `gorm:"type:TEXT;not null;default:''"`
	QaContentTail       string `gorm:"type:TEXT;not null;default:''"`
	QaMobileContentHead string `gorm:"type:TEXT;not null;default:''"`
	QaMobileContentTail string `gorm:"type:TEXT;not null;default:''"`
	Qa1Subj             string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa2Subj             string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa3Subj             string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa4Subj             string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa5Subj             string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa1                 string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa2                 string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa3                 string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa4                 string `gorm:"type:VARCHAR(255);not null;default:''"`
	Qa5                 string `gorm:"type:VARCHAR(255);not null;default:''"`
}

func (*QaConfig) TableName() string {
	return config.Global.DbTablePrefix + TableNameQaConfig
}
