package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zakopokeGo/internal/domain/repository"
)

type pokeAPIRepository struct{}

func NewPokeAPIRepository() repository.PokemonMetadataRepository {
	return &pokeAPIRepository{}
}

func (r *pokeAPIRepository) GetPokemonMetadata(id int) (*repository.PokemonMetadata, error) {
	// 1. 基本データの取得（画像、タイプ）
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch pokemon data: %v", err)
	}
	defer resp.Body.Close()

	var pokeData struct {
		Name    string `json:"name"`
		Sprites struct {
			FrontDefault string `json:"front_default"`
		} `json:"sprites"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pokeData); err != nil {
		return nil, err
	}

	// 2. 日本語名の取得
	speciesURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%d", id)
	speciesResp, err := http.Get(speciesURL)
	japaneseName := pokeData.Name // fallback
	if err == nil && speciesResp.StatusCode == 200 {
		defer speciesResp.Body.Close()
		var speciesData struct {
			Names []struct {
				Language struct {
					Name string `json:"name"`
				} `json:"language"`
				Name string `json:"name"`
			} `json:"names"`
		}
		if err := json.NewDecoder(speciesResp.Body).Decode(&speciesData); err == nil {
			for _, n := range speciesData.Names {
				if n.Language.Name == "ja-Hrkt" {
					japaneseName = n.Name
					break
				}
			}
		}
	}

	// 3. タイプ名の日本語取得
	var types []string
	for _, t := range pokeData.Types {
		typeResp, err := http.Get(t.Type.URL)
		if err != nil || typeResp.StatusCode != 200 {
			types = append(types, t.Type.Name) // fallback
			continue
		}
		defer typeResp.Body.Close()
		var typeData struct {
			Names []struct {
				Language struct {
					Name string `json:"name"`
				} `json:"language"`
				Name string `json:"name"`
			} `json:"names"`
		}
		if err := json.NewDecoder(typeResp.Body).Decode(&typeData); err == nil {
			found := false
			for _, n := range typeData.Names {
				if n.Language.Name == "ja-Hrkt" {
					types = append(types, n.Name)
					found = true
					break
				}
			}
			if !found {
				types = append(types, t.Type.Name)
			}
		} else {
			types = append(types, t.Type.Name)
		}
	}

	return &repository.PokemonMetadata{
		ID:           id,
		Name:         pokeData.Name,
		JapaneseName: japaneseName,
		Image:        pokeData.Sprites.FrontDefault,
		Types:        types,
	}, nil
}
