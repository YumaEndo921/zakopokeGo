package main

import (
	"zakopokeGo/db"
	"zakopokeGo/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// ginフレームワークの初期定義
	r := gin.Default()

	// DB初期化
	db.InitDB()

	// HTMLテンプレートの読み込み
	r.LoadHTMLGlob("templates/*.html")
	// アセット管理の読み込み
	r.Static("/static", "./static")

	// セッションデータをクッキーに保存するためのストアを作成
	store := cookie.NewStore([]byte("secret-gal-key")) // 安全な鍵にしてね〜💘
	// ginにセッションを使うよって教えてる
	r.Use(sessions.Sessions("gal_session", store))

	r.GET("/", handlers.ShowLogin)
	r.POST("/login", handlers.Login)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.Register)
	r.GET("/home", handlers.Home)
	r.GET("/explore", handlers.Explore)
	r.POST("/explore", handlers.Explore)
	r.POST("/catch", handlers.Catch)
	r.POST("/run", handlers.Run)
	r.GET("/logout", handlers.Logout)

	r.GET("/mypokemon", handlers.MyPokemon)

	r.Run(":8080")
}
