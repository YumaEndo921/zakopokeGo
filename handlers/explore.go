package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PokemonData struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func Explore(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	pokeID := rand.Intn(151) + 1 // 1〜151のランダム

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(pokeID))
	if err != nil || resp.StatusCode != 200 {
		c.String(http.StatusInternalServerError, "ポケモン探索失敗しちゃった💦")
		return
	}
	defer resp.Body.Close()

	var data PokemonData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.String(http.StatusInternalServerError, "データ読み込み失敗したかも〜😭")
		return
	}

	// sessionからユーザーID取得（ダミーで UserId 1 を使ってるよ）
	c.HTML(http.StatusOK, "explore_result.html", gin.H{
		"Pokemon": data,
	})
}
