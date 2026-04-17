package external

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
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
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Moves []struct {
			Move struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move"`
		} `json:"moves"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pokeData); err != nil {
		return nil, err
	}

	// Stats の詰め替え
	stats := make(map[string]int)
	for _, s := range pokeData.Stats {
		stats[s.Stat.Name] = s.BaseStat
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
				if n.Language.Name == "ja" {
					japaneseName = n.Name
				}
			}
		}
	}

	// 3. タイプ名の日本語取得
	var types []string
	for _, t := range pokeData.Types {
		typeResp, err := http.Get(t.Type.URL)
		if err != nil || typeResp.StatusCode != 200 {
			types = append(types, t.Type.Name)
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
			var typeName string
			for _, n := range typeData.Names {
				if n.Language.Name == "ja-Hrkt" {
					typeName = n.Name
					break
				}
				if n.Language.Name == "ja" {
					typeName = n.Name
				}
			}
			if typeName != "" {
				types = append(types, typeName)
			} else {
				types = append(types, t.Type.Name)
			}
		} else {
			types = append(types, t.Type.Name)
		}
	}

	// 4. ランダムな攻撃技の取得
	var moves []repository.Move
	// 全ての技を見ると時間がかかるので、ランダムにいくつかピックアップして攻撃技を探す
	rand.Seed(time.Now().UnixNano())
	
	// シャッフルして最初に見つかった「威力がある技」を採用する
	indices := rand.Perm(len(pokeData.Moves))
	count := 0
	for _, idx := range indices {
		if count >= 10 { // 最大10個調べて見つからなければ諦める
			break
		}
		count++

		m := pokeData.Moves[idx]
		moveResp, err := http.Get(m.Move.URL)
		if err != nil || moveResp.StatusCode != 200 {
			continue
		}
		defer moveResp.Body.Close()

		var moveData struct {
			Name  string `json:"name"`
			Power int    `json:"power"`
			Type  struct {
				Name string `json:"name"`
			} `json:"type"`
			Names []struct {
				Language struct {
					Name string `json:"name"`
				} `json:"language"`
				Name string `json:"name"`
			} `json:"names"`
		}
		if err := json.NewDecoder(moveResp.Body).Decode(&moveData); err != nil {
			continue
		}

		if moveData.Power > 0 {
			// 日本語名を探す
			moveJaName := moveData.Name
			for _, n := range moveData.Names {
				if n.Language.Name == "ja-Hrkt" {
					moveJaName = n.Name
					break
				}
				if n.Language.Name == "ja" {
					moveJaName = n.Name
				}
			}

			moves = append(moves, repository.Move{
				Name:  moveJaName,
				Type:  moveData.Type.Name,
				Power: moveData.Power,
			})
			break // 1つ見つかればOK
		}
	}

	// 万が一技が見つからなかった時のフォールバック
	if len(moves) == 0 {
		moves = append(moves, repository.Move{
			Name:  "たいあたり",
			Type:  "normal",
			Power: 40,
		})
	}

	return &repository.PokemonMetadata{
		ID:           id,
		Name:         pokeData.Name,
		JapaneseName: japaneseName,
		Image:        pokeData.Sprites.FrontDefault,
		Types:        types,
		Stats:        stats,
		Moves:        moves,
	}, nil
}
