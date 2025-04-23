package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"zakopokeGo/db"
	"zakopokeGo/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Catch(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	pokemonIDStr := c.PostForm("poke_id")
	pokemonID, err := strconv.Atoi(pokemonIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "home.html", gin.H{"Error": "ポケモンIDが不正だよ〜💦"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	successRate := 70 // 成功率70%
	roll := rand.Intn(100)

	if roll < successRate {
		// 成功時：DB登録
		caught := models.OwnedPokemon{
			UserID:    userID.(uint),
			PokemonNo: pokemonID,
		}
		db.DB.Create(&caught)
		c.HTML(http.StatusOK, "catch_success.html", gin.H{
			"PokemonID": pokemonID,
		})
	} else {
		// 失敗時
		c.HTML(http.StatusOK, "catch_fail.html", gin.H{
			"PokemonID": pokemonID,
		})
	}
}

func Run(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/home")
}
