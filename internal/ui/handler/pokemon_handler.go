package handler

import (
	"net/http"
	"strconv"
	usecase "zakopokeGo/internal/application/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PokemonHandler struct {
	pokemonUC usecase.PokemonUseCase
}

func NewPokemonHandler(pokemonUC usecase.PokemonUseCase) *PokemonHandler {
	return &PokemonHandler{pokemonUC: pokemonUC}
}

func (h *PokemonHandler) Home(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	count, _ := h.pokemonUC.GetCatchCount(userID.(uint))
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Count": count,
	})
}

func (h *PokemonHandler) Explore(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	meta, err := h.pokemonUC.Explore()
	if err != nil {
		c.String(http.StatusInternalServerError, "探索失敗: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "explore_result.html", gin.H{
		"Pokemon": meta,
	})
}

func (h *PokemonHandler) Catch(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	pokeIDStr := c.PostForm("poke_id")
	pokeID, _ := strconv.Atoi(pokeIDStr)

	caught, err := h.pokemonUC.Catch(userID.(uint), pokeID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"Error": "エラーが発生しました"})
		return
	}

	if caught {
		c.HTML(http.StatusOK, "catch_success.html", gin.H{"PokemonID": pokeID})
	} else {
		c.HTML(http.StatusOK, "catch_fail.html", gin.H{"PokemonID": pokeID})
	}
}

func (h *PokemonHandler) Box(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	details, err := h.pokemonUC.GetMyPokemons(userID.(uint))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"Error": "データ取得失敗"})
		return
	}

	c.HTML(http.StatusOK, "box.html", gin.H{
		"PokemonList": details,
	})
}

func (h *PokemonHandler) Run(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/home")
}

func (h *PokemonHandler) Release(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	capturedID, _ := strconv.ParseUint(idStr, 10, 64)

	name, err := h.pokemonUC.ReleasePokemon(userID.(uint), uint(capturedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お別れに失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "バイバイ、" + name + "！ 元気でね！",
	})
}
