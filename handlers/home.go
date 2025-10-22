package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	session := sessions.Default(c)
	UserId := session.Get("user_id")

	if UserId == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "home.html", nil)
}
