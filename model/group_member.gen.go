// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameGroupMember = "group_member"

// GroupMember mapped from table <group_member>
type GroupMember struct {
	GmID       *int32     `gorm:"column:gm_id;type:INTEGER" json:"gm_id"`
	GrID       *string    `gorm:"column:gr_id;type:VARCHAR(10)" json:"gr_id"`
	MbID       *string    `gorm:"column:mb_id;type:VARCHAR(20)" json:"mb_id"`
	GmDatetime *time.Time `gorm:"column:gm_datetime;type:DATETIME" json:"gm_datetime"`
}

// TableName GroupMember's table name
func (*GroupMember) TableName() string {
	return Prefix + TableNameGroupMember
}
