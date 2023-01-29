package model

type Song struct {
	Id       int64  `gorm:"primaryKey;required" json:"id"`
	Album_Id int64  `gorm:"required;ForeignKey:Album" json:"album_id"`
	Title    string `gorm:"required" json:"title"`
	Lyrics   string `gorm:"required" json:"lyrics"`
}
