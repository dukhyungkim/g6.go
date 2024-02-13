// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameWriteFree = "write_free"

// WriteFree mapped from table <write_free>
type WriteFree struct {
	WrID           *int32     `gorm:"column:wr_id;type:INTEGER" json:"wr_id"`
	WrNum          *int32     `gorm:"column:wr_num;type:INTEGER" json:"wr_num"`
	WrReply        *string    `gorm:"column:wr_reply;type:VARCHAR(10)" json:"wr_reply"`
	WrParent       *int32     `gorm:"column:wr_parent;type:INTEGER" json:"wr_parent"`
	WrIsComment    *int32     `gorm:"column:wr_is_comment;type:INTEGER" json:"wr_is_comment"`
	WrComment      *int32     `gorm:"column:wr_comment;type:INTEGER" json:"wr_comment"`
	WrCommentReply *string    `gorm:"column:wr_comment_reply;type:VARCHAR(5)" json:"wr_comment_reply"`
	CaName         *string    `gorm:"column:ca_name;type:VARCHAR(255)" json:"ca_name"`
	WrOption       *string    `gorm:"column:wr_option;type:VARCHAR(40)" json:"wr_option"`
	WrSubject      *string    `gorm:"column:wr_subject;type:VARCHAR(255)" json:"wr_subject"`
	WrContent      *string    `gorm:"column:wr_content;type:TEXT" json:"wr_content"`
	WrSeoTitle     *string    `gorm:"column:wr_seo_title;type:VARCHAR(255)" json:"wr_seo_title"`
	WrLink1        *string    `gorm:"column:wr_link1;type:TEXT" json:"wr_link1"`
	WrLink2        *string    `gorm:"column:wr_link2;type:TEXT" json:"wr_link2"`
	WrLink1Hit     *int32     `gorm:"column:wr_link1_hit;type:INTEGER" json:"wr_link1_hit"`
	WrLink2Hit     *int32     `gorm:"column:wr_link2_hit;type:INTEGER" json:"wr_link2_hit"`
	WrHit          *int32     `gorm:"column:wr_hit;type:INTEGER" json:"wr_hit"`
	WrGood         *int32     `gorm:"column:wr_good;type:INTEGER" json:"wr_good"`
	WrNogood       *int32     `gorm:"column:wr_nogood;type:INTEGER" json:"wr_nogood"`
	MbID           *string    `gorm:"column:mb_id;type:VARCHAR(20)" json:"mb_id"`
	WrPassword     *string    `gorm:"column:wr_password;type:VARCHAR(255)" json:"wr_password"`
	WrName         *string    `gorm:"column:wr_name;type:VARCHAR(255)" json:"wr_name"`
	WrEmail        *string    `gorm:"column:wr_email;type:VARCHAR(255)" json:"wr_email"`
	WrHomepage     *string    `gorm:"column:wr_homepage;type:VARCHAR(255)" json:"wr_homepage"`
	WrDatetime     *time.Time `gorm:"column:wr_datetime;type:DATETIME" json:"wr_datetime"`
	WrFile         *int32     `gorm:"column:wr_file;type:INTEGER" json:"wr_file"`
	WrLast         *string    `gorm:"column:wr_last;type:VARCHAR(30)" json:"wr_last"`
	WrIP           *string    `gorm:"column:wr_ip;type:VARCHAR(255)" json:"wr_ip"`
	WrFacebookUser *string    `gorm:"column:wr_facebook_user;type:VARCHAR(255)" json:"wr_facebook_user"`
	WrTwitterUser  *string    `gorm:"column:wr_twitter_user;type:VARCHAR(255)" json:"wr_twitter_user"`
	Wr1            *string    `gorm:"column:wr_1;type:VARCHAR(255)" json:"wr_1"`
	Wr2            *string    `gorm:"column:wr_2;type:VARCHAR(255)" json:"wr_2"`
	Wr3            *string    `gorm:"column:wr_3;type:VARCHAR(255)" json:"wr_3"`
	Wr4            *string    `gorm:"column:wr_4;type:VARCHAR(255)" json:"wr_4"`
	Wr5            *string    `gorm:"column:wr_5;type:VARCHAR(255)" json:"wr_5"`
	Wr6            *string    `gorm:"column:wr_6;type:VARCHAR(255)" json:"wr_6"`
	Wr7            *string    `gorm:"column:wr_7;type:VARCHAR(255)" json:"wr_7"`
	Wr8            *string    `gorm:"column:wr_8;type:VARCHAR(255)" json:"wr_8"`
	Wr9            *string    `gorm:"column:wr_9;type:VARCHAR(255)" json:"wr_9"`
	Wr10           *string    `gorm:"column:wr_10;type:VARCHAR(255)" json:"wr_10"`
}

// TableName WriteFree's table name
func (*WriteFree) TableName() string {
	return TableNameWriteFree
}
