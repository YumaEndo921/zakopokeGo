package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zakopokeGo/db"
	"zakopokeGo/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MyPokemon(c *gin.Context) {
	session := sessions.Default(c)
	UserId := session.Get("user_id")
	if UserId == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var caught []models.OwnedPokemon
	if err := db.DB.Where("user_id = ?", UserId).Find(&caught).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"Error": "ポケモン取得に失敗したよ〜💦"})
		return
	}

	type PokemonInfo struct {
		ID    int
		Name  string
		Types []string
		Image string
	}

	var pokemonList []PokemonInfo

	for _, p := range caught {
		// ① 日本語名を取得（pokemon-species）
		speciesURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%d", p.PokemonNo)
		speciesResp, err := http.Get(speciesURL)
		if err != nil || speciesResp.StatusCode != 200 {
			continue
		}
		defer speciesResp.Body.Close()

		var speciesData struct {
			Names []struct {
				Language struct {
					Name string `json:"name"`
				} `json:"language"`
				Name string `json:"name"`
			} `json:"names"`
		}
		if err := json.NewDecoder(speciesResp.Body).Decode(&speciesData); err != nil {
			continue
		}

		japaneseName := ""
		for _, name := range speciesData.Names {
			if name.Language.Name == "ja-Hrkt" {
				japaneseName = name.Name
				break
			}
		}
		if japaneseName == "" {
			japaneseName = "名前不明"
		}

		// ② 英語APIで画像とタイプ取得（pokemon）
		pokeURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", p.PokemonNo)
		pokeResp, err := http.Get(pokeURL)
		if err != nil || pokeResp.StatusCode != 200 {
			continue
		}
		defer pokeResp.Body.Close()

		var pokeData struct {
			Sprites struct {
				FrontDefault string `json:"front_default"`
			} `json:"sprites"`
			Types []struct {
				Type struct {
					Name string `json:"name"` // 英語だけどこれでAPI叩く！
					URL  string `json:"url"`
				} `json:"type"`
			} `json:"types"`
		}
		if err := json.NewDecoder(pokeResp.Body).Decode(&pokeData); err != nil {
			continue
		}

		// ③ 各タイプの日本語名を取得する✨
		types := []string{}
		for _, t := range pokeData.Types {
			typeURL := t.Type.URL
			typeResp, err := http.Get(typeURL)
			if err != nil || typeResp.StatusCode != 200 {
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
			if err := json.NewDecoder(typeResp.Body).Decode(&typeData); err != nil {
				continue
			}

			japaneseType := ""
			for _, name := range typeData.Names {
				if name.Language.Name == "ja-Hrkt" {
					japaneseType = name.Name
					break
				}
			}
			if japaneseType != "" {
				types = append(types, japaneseType)
			}
		}

		// 最終的にポケモン情報を追加🌟
		pokemonList = append(pokemonList, PokemonInfo{
			ID:    p.PokemonNo,
			Name:  japaneseName,
			Types: types,
			Image: pokeData.Sprites.FrontDefault,
		})
	}

	c.HTML(http.StatusOK, "mypokemon.html", gin.H{
		"PokemonList": pokemonList,
	})
}
