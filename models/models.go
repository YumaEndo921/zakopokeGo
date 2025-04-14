package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   string `gorm:"unique"`
	Mail     string `gorm:"unique"`
	Password string
}

type OwnedPokemon struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PokemonNo int
}
