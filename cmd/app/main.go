package main

import (
	"log"
	"zakopokeGo/internal/application/usecase"
	"zakopokeGo/internal/infrastructure/external"
	"zakopokeGo/internal/infrastructure/persistence"
	"zakopokeGo/internal/ui/handler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// データベースの初期化
	persistence.InitDB()

	// リポジトリの初期化
	userRepo := persistence.NewUserRepository()
	pokemonRepo := persistence.NewPokemonRepository()
	pokeAPIRepo := external.NewPokeAPIRepository()

	// ユースケースの初期化
	authUC := usecase.NewAuthUseCase(userRepo)
	pokemonUC := usecase.NewPokemonUseCase(pokemonRepo, pokeAPIRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authUC)
	pokemonHandler := handler.NewPokemonHandler(pokemonUC)

	// ルーターの初期化
	r := gin.Default()

	// セッションの設定
	store := cookie.NewStore([]byte("secret-gal-key")) // 以前と同じキーを使用
	r.Use(sessions.Sessions("gal_session", store))

	// テンプレートと静的ファイルの読み込み
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// ルーティング
	r.GET("/", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.ShowRegister)
	r.POST("/register", authHandler.Register)
	r.GET("/logout", authHandler.Logout)

	// 認証が必要なルート
	r.GET("/home", pokemonHandler.Home)
	r.GET("/explore", pokemonHandler.Explore)
	r.POST("/explore", pokemonHandler.Explore) // 元の実装に合わせてGET/POST両方をサポート

	r.POST("/catch", pokemonHandler.Catch)
	r.GET("/box", pokemonHandler.Box)
	r.POST("/box/release/:id", pokemonHandler.Release)
	r.POST("/run", pokemonHandler.Run)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
