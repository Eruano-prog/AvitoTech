package app

import (
	"AvitoTech/internal/config"
	"AvitoTech/internal/controller"
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository/postgres"
	"AvitoTech/internal/service"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	db        *sql.DB
	pool      *dockertest.Pool
	itemsPath = "../entity/items.json"
)

func TestMain(m *testing.M) {
	err := cleanenv.ReadEnv(&config.Configuration)
	if err != nil {
		fmt.Printf("Could not load configuration: %s", err)
		return
	}
	pgCfg := config.Configuration.Database

	pool, err = dockertest.NewPool("")
	if err != nil {
		fmt.Printf("Could not connect to docker: %s", err)
		return
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", pgCfg.Username),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pgCfg.Password),
			fmt.Sprintf("POSTGRES_DB=%s", pgCfg.DBName),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Printf("Could not start resource: %s", err)
		return
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgCfg.Username, pgCfg.Password, hostAndPort, pgCfg.DBName)

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		fmt.Printf("Could not connect to docker: %s", err)
		return
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		balance INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		item TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS history (
		id SERIAL PRIMARY KEY,
		sender_name TEXT NOT NULL,
		receiver_name TEXT NOT NULL,
		amount INTEGER
	);
	`)
	if err != nil {
		fmt.Printf("Could not create table: %s", err)
		return
	}

	code := m.Run()

	if err = db.Close(); err != nil {
		fmt.Printf("Could not close database: %s", err)
	}
	if err = pool.Purge(resource); err != nil {
		fmt.Printf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestApiAuth(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	authRequest := controller.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(authRequest)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResponse controller.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, authResponse.Token)
}

func TestApiInfo(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	authRequest := controller.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(authRequest)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	var authResponse controller.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	require.NoError(t, err)

	req, err = http.NewRequest("GET", server.URL+"/api/info", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+*authResponse.Token)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var infoResponse controller.InfoResponse
	err = json.NewDecoder(resp.Body).Decode(&infoResponse)
	require.NoError(t, err)
	assert.NotNil(t, infoResponse.CoinHistory)
}

func TestApiAuth_InvalidCredentials(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	authRequest := controller.AuthRequest{
		Username: "invaliduser",
		Password: "validpassword",
	}
	body, _ := json.Marshal(authRequest)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	authRequest = controller.AuthRequest{
		Username: "invaliduser",
		Password: "invalidpassword",
	}
	body, _ = json.Marshal(authRequest)
	req, err = http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client = &http.Client{}
	resp, err = client.Do(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestApiBuyItem(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	authRequest := controller.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(authRequest)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	var authResponse controller.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	require.NoError(t, err)

	req, err = http.NewRequest("GET", server.URL+"/api/buy/book", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+*authResponse.Token)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestApiSendCoin(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	authRequest := controller.AuthRequest{
		Username: "anotheruser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(authRequest)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	authRequest = controller.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ = json.Marshal(authRequest)
	req, err = http.NewRequest("POST", server.URL+"/api/auth", bytes.NewBuffer(body))
	require.NoError(t, err)

	client = &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	var authResponse controller.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	require.NoError(t, err)

	sendCoinRequest := controller.SendCoinRequest{
		ToUser: "anotheruser",
		Amount: 10,
	}
	body, _ = json.Marshal(sendCoinRequest)
	req, err = http.NewRequest("POST", server.URL+"/api/sendCoin", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+*authResponse.Token)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestApiInfo_Unauthorized(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Could not sync logger: %s", err)
		}
	}(logger)

	err := entity.LoadItems(logger, itemsPath)
	require.NoError(t, err)

	userRepository := postgres.NewUserRepository(logger, db)
	historyRepository := postgres.NewHistoryRepository(logger, db)
	inventoryRepository := postgres.NewInventoryRepository(logger, db)

	jwtService := service.NewJWTService(logger, "secret")

	authService := service.NewAuthService(logger, userRepository, jwtService)
	infoService := service.NewInfoService(logger, userRepository, historyRepository, inventoryRepository)
	coinService := service.NewCoinService(logger, userRepository, inventoryRepository, historyRepository)

	apiController := controller.NewAPIController(logger, authService, infoService, coinService)

	r := chi.NewRouter()
	apiController.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/api/info", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer wrongToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
