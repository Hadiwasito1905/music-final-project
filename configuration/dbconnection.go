package configuration

import (
	"gorm.io/gorm"
	"music-final-project/model"
	"gorm.io/driver/postgres"

	"os"
)

var DB *gorm.DB

func DatabaseConnect() {
	dsn := os.Getenv("POSTGRES_URL")
	database, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		database.AutoMigrate(&model.Artist{}, &model.Album{}, &model.Song{})
	}

	DB = database
}
