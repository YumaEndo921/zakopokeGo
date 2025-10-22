package handlers

import (
	"log"
	"net/http"
	"strings"
	"zakopokeGo/db"
	"zakopokeGo/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	UserId := strings.TrimSpace(c.PostForm("user_id"))
	password := c.PostForm("password")

	var user models.User
	result := db.DB.Where("user_id = ?", UserId).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("ユーザーIDが見つかりません: %s", UserId)
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "ユーザーIDが間違っています"})
		} else {
			log.Printf("データベースエラー: %v", result.Error)
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "サーバーエラーが発生しました"})
		}
		return
	}

	// パスワードのハッシュを比較
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "メールアドレスまたはパスワードが間違っています"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/home")
}

func ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	mail := c.PostForm("mail")
	UserId := strings.TrimSpace(c.PostForm("user_id"))
	password := c.PostForm("password")

	// 簡易バリデーション
	if UserId == "" || strings.TrimSpace(password) == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"Error": "ユーザーIDとパスワードを入力してください"})
		return
	}
	if len(password) < 4 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"Error": "パスワードは4文字以上にしてください"})
		return
	}

	// パスワードをハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"Error": "登録処理でエラーが発生しました"})
		return
	}

	user := models.User{
		Mail:     mail,
		Password: string(hashed),
		UserId:   UserId,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		// ユニーク制約違反などの判定はDBやドライバ依存なので簡易チェック
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"Error": "そのユーザーIDは既に使われています"})
			return
		}
		log.Printf("failed to create user: %v", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"Error": "登録に失敗しました"})
		return
	}

	// 登録後に自動でログイン（セッションに user_id を保存）
	session := sessions.Default(c)
	// セッション固定攻撃対策としていったんクリア
	session.Clear()
	session.Set("user_id", user.Id)
	if err := session.Save(); err != nil {
		log.Printf("failed to save session for user %d: %v", user.Id, err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"Error": "セッションの開始に失敗しました"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/home")
}
