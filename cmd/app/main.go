package main

import (
	"log"
	"zakopokeGo/internal/application/usecase"
	"zakopokeGo/internal/infrastructure/external"
	"zakopokeGo/internal/infrastructure/persistence"
	"zakopokeGo/internal/domain/service"
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

	// サービスの初期化
	battleService := service.NewBattleService()
	expService := service.NewExperienceService()

	// ユースケースの初期化
	authUC := usecase.NewAuthUseCase(userRepo)
	pokemonUC := usecase.NewPokemonUseCase(pokemonRepo, pokeAPIRepo)
	battleUC := usecase.NewBattleUseCase(pokemonRepo, pokeAPIRepo, battleService, expService)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authUC)
	pokemonHandler := handler.NewPokemonHandler(pokemonUC)
	battleHandler := handler.NewBattleHandler(battleUC)

	// ルーターの初期化
	r := gin.Default()

	// ... (セッション, テンプレート設定などはそのまま)
	store := cookie.NewStore([]byte("secret-gal-key"))
	r.Use(sessions.Sessions("gal_session", store))

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// 認証ミドルウェア (既存の PokemonHandler.Home などから userID を取得するために必要)
	// メモ: 現状は簡単のため各ハンドラー内で userID を取得している
	
	r.GET("/", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.ShowRegister)
	r.POST("/register", authHandler.Register)
	r.GET("/logout", authHandler.Logout)

	r.GET("/home", pokemonHandler.Home)
	r.GET("/explore", pokemonHandler.Explore)
	r.POST("/explore", pokemonHandler.Explore)

	r.POST("/catch", pokemonHandler.Catch)
	r.GET("/box", pokemonHandler.Box)
	r.POST("/box/release/:id", pokemonHandler.Release)
	r.POST("/run", pokemonHandler.Run)

	// バトル
	r.GET("/battle", battleHandler.StartBattle)
	r.POST("/battle/turn", battleHandler.ExecuteTurn)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
