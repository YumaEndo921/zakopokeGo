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

func (h *PokemonHandler) MyPokemon(c *gin.Context) {
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

	// 必要に応じてテンプレートの期待に合わせて変換
	// テンプレートは構造体のフィールドを期待している可能性が高い。repository.PokemonMetadataのキーを合わせるか、テンプレートを調整する。
	// repository.PokemonMetadataはName(英語), JapaneseName(日本語), Image, Typesを持つ。
	// 以前のコードはID, Name, Types, Imageを持つ"PokemonList"を使用していた。
	// Metadataは"JapaneseName"を持っているので、これをテンプレート用に"Name"にマッピングするか、テンプレートを更新する。
	// テンプレートは.Name, .Image, .Types, .IDを使用していると仮定する。
	// 表示用にうまくレンダリングできる構造体を渡す。

	type DisplayPokemon struct {
		ID    int
		Name  string
		Image string
		Types []string
	}
	var displayList []DisplayPokemon
	for _, d := range details {
		displayList = append(displayList, DisplayPokemon{
			ID:    d.ID,
			Name:  d.JapaneseName, // メインの名前として日本語名を使用
			Image: d.Image,
			Types: d.Types,
		})
	}

	c.HTML(http.StatusOK, "mypokemon.html", gin.H{
		"PokemonList": displayList,
	})
}

func (h *PokemonHandler) Run(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/home")
}
