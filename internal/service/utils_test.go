package service

import (
	"AvitoTech/internal/entity"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := entity.LoadItems("../entity/items.json")
	if err != nil {
		fmt.Println("Error loading items.")
		return
	}

	code := m.Run()

	os.Exit(code)
}
