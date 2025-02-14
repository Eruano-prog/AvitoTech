package service

import (
	"AvitoTech/internal/entity"
	"fmt"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Error creating logger")
		return
	}
	err = entity.LoadItems(logger, "../entity/items.json")
	if err != nil {
		logger.Error("Error loading items.")
		return
	}

	code := m.Run()

	os.Exit(code)
}
