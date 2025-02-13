package service

import (
	"AvitoTech/internal/service"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestJWT(t *testing.T) {
	logger := zaptest.NewTestingWriter(t)
	jwtService := service.NewJWTService(logger, "Secret")
}
