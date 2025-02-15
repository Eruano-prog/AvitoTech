// Package app
package app

import (
	"AvitoTech/internal/config"
	"AvitoTech/internal/controller"
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository/postgres"
	"AvitoTech/internal/service"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func waitForConnection(logger *zap.Logger, c *sql.DB) error {
	var err error
	for i := 0; i < 10; i++ {
		err = c.Ping()
		if err == nil {
			logger.Info("connected to database")
			return nil
		}
		logger.Warn("waiting for database...", zap.Int("attempt", i+1))
		time.Sleep(1 * time.Second)
	}
	logger.Fatal("failed to ping database after multiple attempts", zap.Error(err))
	return err
}

func setupApp(logger *zap.Logger, db *sql.DB) (*controller.APIController, error) {
	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, config.Configuration.JwtSecret)

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	return apiController, nil
}

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	c := zap.NewProductionConfig()

	c.Level = zap.NewAtomicLevelAt(zap.WarnLevel)

	logger, err := c.Build()
	if err != nil {
		fmt.Printf("cannot create zap logger: %v", err)
		return
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			fmt.Printf("cannot sync zap logger: %v", err)
		}
	}(logger)

	err = cleanenv.ReadEnv(&config.Configuration)
	if err != nil {
		logger.Fatal("cannot load configuration", zap.Error(err))
		return
	}

	err = entity.LoadItems(logger, config.Configuration.ItemsPath)
	if err != nil {
		logger.Fatal("cannot load items", zap.Error(err))
		return
	}

	pgCfg := config.Configuration.Database

	pgAddr := pgCfg.Address
	pgDB := pgCfg.DBName
	pgUser := pgCfg.Username
	pgPass := pgCfg.Password

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPass, pgAddr, pgDB)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return
	}
	db.SetMaxOpenConns(700)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = waitForConnection(logger, db)

	apiController, err := setupApp(logger, db)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Fatal("failed to close database connection", zap.Error(err))
		}
	}(db)
	if err != nil {
		logger.Fatal("failed to setup app", zap.Error(err))
		return
	}

	r := chi.NewRouter()

	apiController.Register(r)

	server := http.Server{
		Addr:    config.Configuration.Server.RESTAddr,
		Handler: r,
	}

	logger.Info("starting server", zap.String("addr", config.Configuration.Server.RESTAddr))
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("cannot start server", zap.Error(err))
	}
}
