package handler

import (
	"net/http"
	usecase "zakopokeGo/internal/application/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	userID := c.PostForm("user_id")
	password := c.PostForm("password")

	user, err := h.authUC.Login(userID, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "ログイン失敗: " + err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/home")
}

func (h *AuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	userID := c.PostForm("user_id")
	mail := c.PostForm("mail")
	password := c.PostForm("password")

	if err := h.authUC.Register(userID, mail, password); err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"Error": "登録失敗: " + err.Error()})
		return
	}

	// 自動ログイン (このロジックも改善できるが、現状シンプルに保つ)
	// とりあえずログインまたはホームへリダイレクト。IDを取得できればホームへ、そうでなければログインへ。
	// RegisterはUserオブジェクトを返さないので、ログインロジックへリダイレクトするか、手順を複製するか。
	// 安全のため、ログインページへリダイレクトする。
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
