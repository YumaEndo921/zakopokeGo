package handlers

import (
	"net/http"
	"zakopokeGo/db"
	"zakopokeGo/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	mail := c.PostForm("mail")
	password := c.PostForm("password")

	var user models.User
	result := db.DB.Where("mail = ? AND password = ?", mail, password).First(&user)
	if result.Error != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "ログイン失敗💦"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/home")
}

func ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	mail := c.PostForm("mail")
	password := c.PostForm("password")
	userID := c.PostForm("user_id")

	var exists models.User
	if err := db.DB.Where("mail = ?", mail).First(&exists).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"Error": "すでに登録されてるよ〜🥺"})
		return
	}

	user := models.User{
		Mail:     mail,
		Password: password, // 本番ならハッシュしようね！
		UserID:   userID,
	}
	db.DB.Create(&user)

	c.Redirect(http.StatusSeeOther, "/home")
}
