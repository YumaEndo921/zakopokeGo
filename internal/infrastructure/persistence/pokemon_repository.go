package persistence

import (
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"
)

type pokemonRepository struct{}

func NewPokemonRepository() repository.PokemonRepository {
	return &pokemonRepository{}
}

func (r *pokemonRepository) Create(pokemon *model.Pokemon) error {
	dbPokemon := FromDomainPokemon(pokemon)
	result := DB.Create(dbPokemon)
	if result.Error != nil {
		return result.Error
	}
	pokemon.ID = dbPokemon.ID
	return nil
}

func (r *pokemonRepository) FindByUserID(userID uint) ([]*model.Pokemon, error) {
	var dbPokemons []Pokemon
	result := DB.Where("user_id = ?", userID).Find(&dbPokemons)
	if result.Error != nil {
		return nil, result.Error
	}

	var pokemons []*model.Pokemon
	for _, p := range dbPokemons {
		pokemons = append(pokemons, p.ToDomain())
	}
	return pokemons, nil
}

func (r *pokemonRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	result := DB.Model(&Pokemon{}).Where("user_id = ?", userID).Count(&count)
	return count, result.Error
}

func (r *pokemonRepository) FindByID(id uint) (*model.Pokemon, error) {
	var dbPokemon Pokemon
	result := DB.First(&dbPokemon, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return dbPokemon.ToDomain(), nil
}

func (r *pokemonRepository) Delete(id uint) error {
	result := DB.Delete(&Pokemon{}, id)
	return result.Error
}

func (r *pokemonRepository) Save(pokemon *model.Pokemon) error {
	dbPokemon := FromDomainPokemon(pokemon)
	result := DB.Save(dbPokemon)
	return result.Error
}
