package entity

import (
	"encoding/json"
	"fmt"
	"os"
)

var Items map[string]int

type Item struct {
	Id      int
	OwnerId int
	Title   string
}

func LoadItems(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Items)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return err
	}
	return nil
}
