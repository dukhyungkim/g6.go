package model

import (
	"time"
)

const TableNameVisit = "visit"

// Visit mapped from table <visit>
type Visit struct {
	ViID      uint      `gorm:"type:INTEGER;primaryKey;autoincrement"`
	ViIP      string    `gorm:"type:VARCHAR(100);not null;default:''"`
	ViDate    time.Time `gorm:"type:DATE;not null;default:''"`
	ViTime    time.Time `gorm:"type:TIME;not null;default:''"`
	ViReferer string    `gorm:"type:TEXT;not null;default:''"`
	ViAgent   string    `gorm:"type:VARCHAR(200);not null;default:''"`
	ViBrowser string    `gorm:"type:VARCHAR(255);not null;default:''"`
	ViOs      string    `gorm:"type:VARCHAR(255);not null;default:''"`
	ViDevice  string    `gorm:"type:VARCHAR(255);not null;default:''"`
}

// TableName Visit's table name
func (*Visit) TableName() string {
	return Prefix + TableNameVisit
}
