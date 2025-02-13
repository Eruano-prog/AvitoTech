package app

import (
	"AvitoTech/internal/config"
	"AvitoTech/internal/controller"
	"AvitoTech/internal/repository/postgres"
	"AvitoTech/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("cannot create zap logger: %v", err)
		return
	}
	defer logger.Sync()

	err = config.LoadConfiguration()
	if err != nil {
		logger.Fatal("cannot load configuration", zap.Error(err))
		return
	}

	pgCfg := config.Configuration.Database.Postgres

	pgAddr := pgCfg.Address
	pdDb := pgCfg.DBName
	pgUser := pgCfg.Username
	pgPass := pgCfg.Password

	userRepository, err := postgres.NewUserRepository(logger, pgAddr, pgUser, pgPass, pdDb)
	if err != nil {
		logger.Fatal("cannot create user repository", zap.Error(err))
		return
	}
	historyRepository, err := postgres.NewHistoryRepository(logger, pgAddr, pgUser, pgPass, pdDb)
	if err != nil {
		logger.Fatal("cannot create history repository", zap.Error(err))
		return
	}
	inventoryRepository, err := postgres.NewInventoryRepository(logger, pgAddr, pgUser, pgPass, pdDb)
	if err != nil {
		logger.Fatal("cannot create inventory repository", zap.Error(err))
		return
	}

	jwtService := service.NewJWTService(logger, config.Configuration.JwtSecret)

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository)

	apiController := controller.NewApiController(logger, authService, infoService, coinService)

	r := chi.NewRouter()

	apiController.Register(r)

	server := http.Server{
		Addr:    config.Configuration.Server.RESTAddr,
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("cannot start server", zap.Error(err))
	}
}
