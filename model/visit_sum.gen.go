// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVisitSum = "visit_sum"

// VisitSum mapped from table <visit_sum>
type VisitSum struct {
	VsDate  *time.Time `gorm:"column:vs_date;type:DATE" json:"vs_date"`
	VsCount *int32     `gorm:"column:vs_count;type:INTEGER" json:"vs_count"`
}

// TableName VisitSum's table name
func (*VisitSum) TableName() string {
	return Prefix + TableNameVisitSum
}
