package main

import (
	"zakopokeGo/db"
	"zakopokeGo/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.InitDB()

	r.LoadHTMLGlob("templates/*.html")

	store := cookie.NewStore([]byte("secret-gal-key")) // 安全な鍵にしてね〜💘
	r.Use(sessions.Sessions("gal_session", store))

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.Login)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.Register)
	r.GET("/home", handlers.Home)
	r.GET("/explore", handlers.Explore)
	r.POST("/explore", handlers.Explore)
	r.POST("/catch", handlers.Catch)
	r.POST("/run", handlers.Run)
	r.GET("/mypokemon", handlers.MyPokemon)

	r.Run(":8080")
}
