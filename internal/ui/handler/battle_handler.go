package handler

import (
	"net/http"
	"strconv"
	"zakopokeGo/internal/application/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BattleHandler struct {
	battleUseCase usecase.BattleUseCase
}

func NewBattleHandler(battleUseCase usecase.BattleUseCase) *BattleHandler {
	return &BattleHandler{
		battleUseCase: battleUseCase,
	}
}

func (h *BattleHandler) StartBattle(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	myPokemonIDStr := c.Query("my_pokemon_id")
	myPokemonID, _ := strconv.ParseUint(myPokemonIDStr, 10, 64)

	result, err := h.battleUseCase.StartBattle(userID.(uint), uint(myPokemonID))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "battle.html", gin.H{
		"Result": result,
	})
}

func (h *BattleHandler) ExecuteTurn(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	myPokemonID, _ := strconv.ParseUint(c.PostForm("my_pokemon_id"), 10, 64)
	action := c.PostForm("action")
	myName := c.PostForm("my_name")
	enemyName := c.PostForm("enemy_name")

	// 敵の状態を隠しフィールドから復元
	enemyNo, _ := strconv.Atoi(c.PostForm("enemy_no"))
	enemyHP, _ := strconv.Atoi(c.PostForm("enemy_hp"))
	enemyMaxHP, _ := strconv.Atoi(c.PostForm("enemy_max_hp"))
	enemyAtk, _ := strconv.Atoi(c.PostForm("enemy_atk"))
	enemyDef, _ := strconv.Atoi(c.PostForm("enemy_def"))
	enemyMoveName := c.PostForm("enemy_move_name")
	enemyMoveType := c.PostForm("enemy_move_type")
	enemyMovePower, _ := strconv.Atoi(c.PostForm("enemy_move_power"))
	enemyType1 := c.PostForm("enemy_type1")
	enemyType2 := c.PostForm("enemy_type2")

	result, err := h.battleUseCase.ExecuteTurn(
		userID.(uint), uint(myPokemonID), action, myName, enemyName,
		enemyNo, enemyHP, enemyMaxHP, enemyAtk, enemyDef,
		enemyMoveName, enemyMoveType, enemyMovePower, enemyType1, enemyType2,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "battle.html", gin.H{
		"Result": result,
	})
}
