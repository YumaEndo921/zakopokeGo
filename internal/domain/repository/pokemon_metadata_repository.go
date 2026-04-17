package repository

type Move struct {
	Name  string
	Type  string
	Power int
}

type PokemonMetadata struct {
	ID           int
	Name         string
	JapaneseName string
	Image        string
	Types        []string
	Stats        map[string]int
	Moves        []Move
}

type PokemonMetadataRepository interface {
	GetPokemonMetadata(pokemonID int) (*PokemonMetadata, error)
}
