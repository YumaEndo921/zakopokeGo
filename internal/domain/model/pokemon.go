package model

import "time"

type Pokemon struct {
	ID        uint
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
	CapturedAt time.Time
}

func NewPokemon(userID uint, pokemonNo int) *Pokemon {
	return &Pokemon{
		UserID:    userID,
		PokemonNo: pokemonNo,
		Level:     1,
		Exp:       0,
	}
}
