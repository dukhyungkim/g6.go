package model

import (
	"time"
)

type WriteTable struct {
	tableName string `gorm:"-"`

	WrID           uint      `gorm:"type:INTEGER;primaryKey;not null"`
	WrNum          int       `gorm:"type:INTEGER;not null;default:0;index:idx_wr_num_reply"`
	WrReply        string    `gorm:"type:VARCHAR(10);not null;default:'';index:idx_wr_num_reply"`
	WrParent       int       `gorm:"type:INTEGER;not null;default:0"`
	WrIsComment    int       `gorm:"type:INTEGER;not null;default:0;index:idx_wr_is_comment"`
	WrComment      int       `gorm:"type:INTEGER;not null;default:0"`
	WrCommentReply string    `gorm:"type:VARCHAR(5);not null;default:''"`
	CaName         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrOption       string    `gorm:"type:VARCHAR(40);not null;default:''"`
	WrSubject      string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrContent      string    `gorm:"type:TEXT;not null;default:''"`
	WrSeoTitle     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrLink1        string    `gorm:"type:TEXT;not null;default:''"`
	WrLink2        string    `gorm:"type:TEXT;not null;default:''"`
	WrLink1Hit     int       `gorm:"type:INTEGER;not null;default:0"`
	WrLink2Hit     int       `gorm:"type:INTEGER;not null;default:0"`
	WrHit          int       `gorm:"type:INTEGER;not null;default:0"`
	WrGood         int       `gorm:"type:INTEGER;not null;default:0"`
	WrNogood       int       `gorm:"type:INTEGER;not null;default:0"`
	MbID           string    `gorm:"type:VARCHAR(20);not null;default:''"`
	WrPassword     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrName         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrEmail        string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrHomepage     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrDatetime     time.Time `gorm:"type:DATETIME;not null;default:''"`
	WrFile         int       `gorm:"type:INTEGER;not null;default:0"`
	WrLast         string    `gorm:"type:VARCHAR(30);not null;default:''"`
	WrIP           string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrFacebookUser string    `gorm:"type:VARCHAR(255);not null;default:''"`
	WrTwitterUser  string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr1            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr2            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr3            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr4            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr5            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr6            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr7            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr8            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr9            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Wr10           string    `gorm:"type:VARCHAR(255);not null;default:''"`
}

func NewWriteTable(tableName string) *WriteTable {
	return &WriteTable{tableName: tableName}
}

func (t *WriteTable) TableName() string {
	return Prefix + WriteTablePrefix + t.tableName
}
