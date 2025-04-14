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
	userID := session.Get("user_id").(uint)
	pokeID, _ := strconv.Atoi(c.PostForm("poke_id"))

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < 60 { // 60%で捕獲成功だYO〜★
		db.DB.Create(&models.OwnedPokemon{
			UserID:    uint(userID),
			PokemonNo: pokeID,
		})

		c.HTML(http.StatusOK, "catch_result.html", gin.H{
			"Success":   true,
			"PokemonID": pokeID,
		})
	} else {
		c.HTML(http.StatusOK, "catch_result.html", gin.H{
			"Success":   false,
			"PokemonID": pokeID,
		})
	}
}

func Run(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/home")
}
