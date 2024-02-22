package model

const TableNameFaqMaster = "faq_master"

// FaqMaster mapped from table <faq_master>
type FaqMaster struct {
	FmID             uint   `gorm:"type:INTEGER;primaryKey;autoIncrement"`
	FmSubject        string `gorm:"type:VARCHAR(255);not null;default:''"`
	FmHeadHTML       string `gorm:"type:TEXT;not null;default:''"`
	FmTailHTML       string `gorm:"type:TEXT;not null;default:''"`
	FmMobileHeadHTML string `gorm:"type:TEXT;not null;default:''"`
	FmMobileTailHTML string `gorm:"type:TEXT;not null;default:''"`
	FmOrder          int    `gorm:"type:INTEGER;not null;default:0"`
}

// TableName FaqMaster's table name
func (*FaqMaster) TableName() string {
	return Prefix + TableNameFaqMaster
}
