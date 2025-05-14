package controller

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/internal/api/controller/auth"
	"github.com/AlexJudin/wallet_java_code/internal/api/controller/register"
	"github.com/AlexJudin/wallet_java_code/internal/api/controller/wallet"
	"github.com/AlexJudin/wallet_java_code/internal/api/middleware"
	"github.com/AlexJudin/wallet_java_code/internal/cache"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
	"github.com/AlexJudin/wallet_java_code/internal/service"
	"github.com/AlexJudin/wallet_java_code/internal/usecases"
)

func AddRoutes(config *config.Сonfig, db *gorm.DB, redisClient *redis.Client, r *chi.Mux) {
	// init services
	authService := service.NewAuthService(config)

	// init repository
	repoWallet := repository.NewWalletRepo(db)
	repoUser := repository.NewUserRepo(db)

	// init cache
	cache := cache.NewCacheClientRepo(redisClient)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repoWallet, cache)
	walletHandler := wallet.NewWalletHandler(walletUC)

	registerUC := usecases.NewRegisterUsecase(repoUser, authService)
	registerHandler := register.NewRegisterHandler(registerUC)

	authUC := usecases.NewAuthUsecase(repoUser, authService)
	authHandler := auth.NewAuthHandler(authUC)

	// init middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	r.Post("/register", registerHandler.RegisterUser)
	r.Post("/auth", authHandler.AuthorizationUser)
	r.Post("/refresh-token", authHandler.RefreshToken)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5000, time.Second))
		r.Use(authMiddleware.CheckToken)
		r.Post("/api/v1/wallet", walletHandler.CreateOperation)
		r.Get("/api/v1/wallets/", walletHandler.GetWalletBalanceByUUID)
	})
}
