package usecase

import (
	"math/rand"
	"time"
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"
)

type CapturedPokemonDetail struct {
	ID           uint
	PokemonNo    int
	Name         string
	JapaneseName string
	Image        string
	Types        []string
}

type PokemonUseCase interface {
	Explore() (*repository.PokemonMetadata, error)
	Catch(userID uint, pokemonNo int) (bool, error)
	GetMyPokemons(userID uint) ([]*CapturedPokemonDetail, error)
	GetCatchCount(userID uint) (int64, error)
	ReleasePokemon(userID uint, capturedID uint) (string, error)
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

func (u *pokemonUseCase) GetMyPokemons(userID uint) ([]*CapturedPokemonDetail, error) {
	pokemons, err := u.pokemonRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var details []*CapturedPokemonDetail
	for _, p := range pokemons {
		meta, err := u.metadataRepo.GetPokemonMetadata(p.PokemonNo)
		if err != nil {
			continue
		}
		details = append(details, &CapturedPokemonDetail{
			ID:           p.ID,
			PokemonNo:    p.PokemonNo,
			Name:         meta.Name,
			JapaneseName: meta.JapaneseName,
			Image:        meta.Image,
			Types:        meta.Types,
		})
	}
	return details, nil
}

func (u *pokemonUseCase) GetCatchCount(userID uint) (int64, error) {
	return u.pokemonRepo.CountByUserID(userID)
}

func (u *pokemonUseCase) ReleasePokemon(userID uint, capturedID uint) (string, error) {
	// 1. 所有権の確認
	pokemon, err := u.pokemonRepo.FindByID(capturedID)
	if err != nil {
		return "", err
	}
	if pokemon.UserID != userID {
		return "", model.ErrUnauthorized // 適当なエラー
	}

	// 2. ポケモン名の取得（演出用）
	meta, err := u.metadataRepo.GetPokemonMetadata(pokemon.PokemonNo)
	name := "ポケモン"
	if err == nil {
		name = meta.JapaneseName
	}

	// 3. 削除
	err = u.pokemonRepo.Delete(capturedID)
	if err != nil {
		return "", err
	}

	return name, nil
}
