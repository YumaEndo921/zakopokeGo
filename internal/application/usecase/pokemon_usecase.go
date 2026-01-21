package usecase

import (
	"math/rand"
	"time"
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"
)

type PokemonUseCase interface {
	Explore() (*repository.PokemonMetadata, error)
	Catch(userID uint, pokemonNo int) (bool, error)
	GetMyPokemons(userID uint) ([]*repository.PokemonMetadata, error)
	GetCatchCount(userID uint) (int64, error)
}

type pokemonUseCase struct {
	pokemonRepo  repository.PokemonRepository
	metadataRepo repository.PokemonMetadataRepository
}

func NewPokemonUseCase(pokemonRepo repository.PokemonRepository, metadataRepo repository.PokemonMetadataRepository) PokemonUseCase {
	return &pokemonUseCase{
		pokemonRepo:  pokemonRepo,
		metadataRepo: metadataRepo,
	}
}

func (u *pokemonUseCase) Explore() (*repository.PokemonMetadata, error) {
	rand.Seed(time.Now().UnixNano())
	pokemonID := rand.Intn(151) + 1
	return u.metadataRepo.GetPokemonMetadata(pokemonID)
}

func (u *pokemonUseCase) Catch(userID uint, pokemonNo int) (bool, error) {
	rand.Seed(time.Now().UnixNano())
	success := rand.Intn(2) == 0 // 50% の確率

	if success {
		newPokemon := model.NewPokemon(userID, pokemonNo)
		err := u.pokemonRepo.Create(newPokemon)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (u *pokemonUseCase) GetMyPokemons(userID uint) ([]*repository.PokemonMetadata, error) {
	pokemons, err := u.pokemonRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var details []*repository.PokemonMetadata
	for _, p := range pokemons {
		meta, err := u.metadataRepo.GetPokemonMetadata(p.PokemonNo)
		if err != nil {
			continue // メタデータの取得に失敗した場合はスキップ、またはエラーハンドリング
		}
		details = append(details, meta)
	}
	return details, nil
}

func (u *pokemonUseCase) GetCatchCount(userID uint) (int64, error) {
	return u.pokemonRepo.CountByUserID(userID)
}
