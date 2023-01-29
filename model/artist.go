package model

type Artist struct {
	Id     int64   `gorm:"primaryKey;required" json:"id"`
	Name   string  `gorm:"required" json:"name"`
	Albums []Album `gorm:"ForeignKey:Artist_Id"`
}
