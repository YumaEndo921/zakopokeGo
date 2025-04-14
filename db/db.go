package db

import (
	"zakopokeGo/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("pokeapp.db"), &gorm.Config{})
	if err != nil {
		panic("DB開けなかった〜😫")
	}

	DB.AutoMigrate(&models.User{}, &models.OwnedPokemon{})
}
