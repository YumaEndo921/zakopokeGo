package model

type Pokemon struct {
	ID        uint
	UserID    uint
	PokemonNo int
}

func NewPokemon(userID uint, pokemonNo int) *Pokemon {
	return &Pokemon{
		UserID:    userID,
		PokemonNo: pokemonNo,
	}
}
