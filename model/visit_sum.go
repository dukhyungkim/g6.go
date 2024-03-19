package model

const TableNameVisitSum = "visit_sum"

// VisitSum mapped from table <visit_sum>
type VisitSum struct {
	VsDate  string `gorm:"type:DATE;primaryKey;not null;default:''" json:"vs_date"`
	VsCount int64  `gorm:"type:INTEGER;not null;default:0" json:"vs_count"`
}

// TableName VisitSum's table name
func (*VisitSum) TableName() string {
	return Prefix + TableNameVisitSum
}
