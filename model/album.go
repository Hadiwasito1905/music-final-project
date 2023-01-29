package model

type Album struct {
	Id        int64  `gorm:"primaryKey;required" json:"id"`
	Artist_Id int64  `gorm:"required;ForeignKey:Artist" json:"artist_id"`
	Title     string `gorm:"required" json:"title"`
	Price     int64  `gorm:"required" json:"price"`
	Songs     []Song `gorm:"ForeignKey:Album_Id"`
}
