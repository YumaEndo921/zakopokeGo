package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear() // セッションの中身をぜ〜んぶ消す💥
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
