package persistence

import (
	"log"
	"zakopokeGo/internal/domain/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBモデル（ドメインモデルとは別）
type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   string `gorm:"unique"`
	Mail     string `gorm:"unique"`
	Password string
}

type Pokemon struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PokemonNo int
	Level     int
	Exp       int
	MaxHP     int
	CurrentHP int
	Attack    int
	Defense   int
	MoveName  string
	MoveType  string
	MovePower int
	Type1     string
	Type2     string
}

func (u *User) ToDomain() *model.User {
	return &model.User{
		ID:       u.ID,
		UserID:   u.UserID,
		Mail:     u.Mail,
		Password: u.Password,
	}
}

func FromDomainUser(u *model.User) *User {
	return &User{
		ID:       u.ID,
		UserID:   u.UserID,
		Mail:     u.Mail,
		Password: u.Password,
	}
}

func (p *Pokemon) ToDomain() *model.Pokemon {
	return &model.Pokemon{
		ID:        p.ID,
		UserID:    p.UserID,
		PokemonNo: p.PokemonNo,
		Level:     p.Level,
		Exp:       p.Exp,
		MaxHP:     p.MaxHP,
		CurrentHP: p.CurrentHP,
		Attack:    p.Attack,
		Defense:   p.Defense,
		MoveName:  p.MoveName,
		MoveType:  p.MoveType,
		MovePower: p.MovePower,
		Type1:     p.Type1,
		Type2:     p.Type2,
	}
}

func FromDomainPokemon(p *model.Pokemon) *Pokemon {
	return &Pokemon{
		ID:        p.ID,
		UserID:    p.UserID,
		PokemonNo: p.PokemonNo,
		Level:     p.Level,
		Exp:       p.Exp,
		MaxHP:     p.MaxHP,
		CurrentHP: p.CurrentHP,
		Attack:    p.Attack,
		Defense:   p.Defense,
		MoveName:  p.MoveName,
		MoveType:  p.MoveType,
		MovePower: p.MovePower,
		Type1:     p.Type1,
		Type2:     p.Type2,
	}
}

// DBグローバル変数（DIすることも可能だが、現状はシンプルに保つ）
var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("pokeapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	DB.AutoMigrate(&User{}, &Pokemon{})
}
