package app

import (
	"AvitoTech/internal/config"
	"AvitoTech/internal/controller"
	"AvitoTech/internal/repository/postgres"
	"AvitoTech/internal/service"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	pgDb := pgCfg.DBName
	pgUser := pgCfg.Username
	pgPass := pgCfg.Password

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPass, pgAddr, pgDb)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return
	}

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, config.Configuration.JwtSecret)

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

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
