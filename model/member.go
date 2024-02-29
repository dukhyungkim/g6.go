package model

import (
	"time"
)

const TableNameMember = "member"

// Member mapped from table <member>
type Member struct {
	MbNo            uint      `gorm:"type:INTEGER;primaryKey;not null"`
	MbID            string    `gorm:"type:VARCHAR(20);unique;not null;default:''"`
	MbPassword      string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbName          string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbNick          string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbNickDate      time.Time `gorm:"type:DATE;not null;default:'0001-01-01'"`
	MbEmail         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbHomepage      string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbLevel         int       `gorm:"type:INTEGER;not null;default:0"`
	MbSex           string    `gorm:"type:VARCHAR(1);not null;default:''"`
	MbBirth         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbTel           string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbHp            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbCertify       string    `gorm:"type:VARCHAR(20);not null;default:''"`
	MbAdult         int       `gorm:"type:INTEGER;not null;default:0"`
	MbDupinfo       string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbZip1          string    `gorm:"type:VARCHAR(3);not null;default:''"`
	MbZip2          string    `gorm:"type:VARCHAR(3);not null;default:''"`
	MbAddr1         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbAddr2         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbAddr3         string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbAddrJibeon    string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbSignature     string    `gorm:"type:TEXT;not null;default:''"`
	MbRecommend     string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbPoint         int       `gorm:"type:INTEGER;not null;default:0"`
	MbTodayLogin    time.Time `gorm:"type:DATETIME;not null;default:'0001-01-01'"`
	MbLoginIP       string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbDatetime      time.Time `gorm:"type:DATETIME;not null;default:'0001-01-01'"`
	MbIP            string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbLeaveDate     string    `gorm:"type:VARCHAR(8);not null;default:''"`
	MbInterceptDate string    `gorm:"type:VARCHAR(8);not null;default:''"`
	MbEmailCertify  time.Time `gorm:"type:DATETIME;not null;default:'0001-01-01'"`
	MbEmailCertify2 string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbMemo          string    `gorm:"type:TEXT;not null;default:''"`
	MbLostCertify   string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbMailling      int       `gorm:"type:INTEGER;not null;default:0"`
	MbSms           int       `gorm:"type:INTEGER;not null;default:0"`
	MbOpen          int       `gorm:"type:INTEGER;not null;default:0"`
	MbOpenDate      time.Time `gorm:"type:DATE;not null;default:'0001-01-01'"`
	MbProfile       string    `gorm:"type:TEXT;not null;default:''"`
	MbMemoCall      string    `gorm:"type:VARCHAR(255);not null;default:''"`
	MbMemoCnt       int       `gorm:"type:INTEGER;not null;default:0"`
	MbScrapCnt      int       `gorm:"type:INTEGER;not null;default:0"`
	Mb1             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb2             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb3             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb4             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb5             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb6             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb7             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb8             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb9             string    `gorm:"type:VARCHAR(255);not null;default:''"`
	Mb10            string    `gorm:"type:VARCHAR(255);not null;default:''"`
}

// TableName Member's table name
func (*Member) TableName() string {
	return Prefix + TableNameMember
}

func (m *Member) IsInterceptOrLeave() bool {
	if m.MbID == "" {
		return false
	}

	return m.MbLeaveDate != "" || m.MbInterceptDate != ""
}
