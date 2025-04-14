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
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var caught []models.OwnedPokemon
	if err := db.DB.Where("user_id = ?", userID).Find(&caught).Error; err != nil {
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
		url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", p.PokemonNo)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			continue
		}

		// 名前と画像とタイプを抽出✨
		name := result["name"].(string)
		image := result["sprites"].(map[string]interface{})["front_default"].(string)

		// タイプの抽出（複数ある場合もある！）
		types := []string{}
		for _, t := range result["types"].([]interface{}) {
			typeInfo := t.(map[string]interface{})["type"].(map[string]interface{})["name"].(string)
			types = append(types, typeInfo)
		}

		pokemonList = append(pokemonList, PokemonInfo{
			ID:    p.PokemonNo,
			Name:  name,
			Types: types,
			Image: image,
		})
	}

	c.HTML(http.StatusOK, "mypokemon.html", gin.H{
		"PokemonList": pokemonList,
	})
}
