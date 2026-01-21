package repository

type PokemonMetadata struct {
	ID           int
	Name         string
	JapaneseName string
	Image        string
	Types        []string
}

type PokemonMetadataRepository interface {
	GetPokemonMetadata(pokemonID int) (*PokemonMetadata, error)
}
