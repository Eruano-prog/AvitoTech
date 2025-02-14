// Package entity
package entity

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
)

var Items map[string]int

type Item struct {
	ID      int
	OwnerID int
	Title   string
}

func LoadItems(l *zap.Logger, path string) error {
	file, err := os.Open(path)
	if err != nil {
		l.Error("Error opening file:", zap.Error(err))
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			l.Error("Error closing file", zap.Error(err))
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Items)
	if err != nil {
		l.Error("Error decoding JSON:", zap.Error(err))
		return err
	}
	return nil
}
