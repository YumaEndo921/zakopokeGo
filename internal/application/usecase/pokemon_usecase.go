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
	Level        int
	CapturedAt   time.Time
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
		meta, err := u.metadataRepo.GetPokemonMetadata(pokemonNo)
		if err != nil {
			return false, err
		}

		newPokemon := model.NewPokemon(userID, pokemonNo)
		// ステータスの設定 (PokeAPIの形式に合わせてマッピング)
		// hp, attack, defense
		newPokemon.MaxHP = meta.Stats["hp"]
		newPokemon.CurrentHP = newPokemon.MaxHP
		newPokemon.Attack = meta.Stats["attack"]
		newPokemon.Defense = meta.Stats["defense"]

		// タイプの共通化
		if len(meta.Types) > 0 {
			newPokemon.Type1 = meta.Types[0]
			if len(meta.Types) > 1 {
				newPokemon.Type2 = meta.Types[1]
			}
		}

		// 技の設定 (metadataRepo でランダムに1つ選ばれている想定)
		if len(meta.Moves) > 0 {
			newPokemon.MoveName = meta.Moves[0].Name
			newPokemon.MoveType = meta.Moves[0].Type
			newPokemon.MovePower = meta.Moves[0].Power
		}

		err = u.pokemonRepo.Create(newPokemon)
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
			Level:        p.Level,
			CapturedAt:   p.CapturedAt,
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
