package models

type User struct {
	Id       uint   `gorm:"primaryKey"`
	UserId   string `gorm:"unique"`
	Mail     string `gorm:"unique"`
	Password string
}

type OwnedPokemon struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint
	PokemonNo int
}
