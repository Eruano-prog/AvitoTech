package postgres

import (
	"AvitoTech/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestInsertItem(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewInventoryRepository(logger, db)

	item, err := repo.InsertItem(1, "item1")
	assert.NoError(t, err)
	defer func(repo repository.InventoryRepository, id int) {
		err = repo.DeleteItem(id)
		if err != nil {
			logger.Error("error deleting item", zap.Error(err))
		}
	}(repo, item.ID)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM inventory WHERE owner_id = $1 AND item = $2", 1, "item1").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetUsersInventory(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewInventoryRepository(logger, db)

	item1, err := repo.InsertItem(1, "item1")
	assert.NoError(t, err)
	defer func(repo repository.InventoryRepository, id int) {
		err = repo.DeleteItem(id)
		if err != nil {
			logger.Error("error deleting item", zap.Error(err))
		}
	}(repo, item1.ID)

	item2, err := repo.InsertItem(1, "item1")
	assert.NoError(t, err)
	defer func(repo repository.InventoryRepository, id int) {
		err = repo.DeleteItem(id)
		if err != nil {
			logger.Error("error deleting item", zap.Error(err))
		}
	}(repo, item2.ID)

	item3, err := repo.InsertItem(1, "item2")
	assert.NoError(t, err)
	defer func(repo repository.InventoryRepository, id int) {
		err = repo.DeleteItem(id)
		if err != nil {
			logger.Error("error deleting item", zap.Error(err))
		}
	}(repo, item3.ID)

	inventory, err := repo.GetUsersInventory(1)
	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, 2, inventory["item1"])
	assert.Equal(t, 1, inventory["item2"])
}

func TestDeleteItem(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewInventoryRepository(logger, db)

	item, err := repo.InsertItem(1, "item1")
	assert.NoError(t, err)

	err = repo.DeleteItem(item.ID)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM inventory WHERE id = $1", item.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetUsersInventoryEmpty(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewInventoryRepository(logger, db)

	inventory, err := repo.GetUsersInventory(999)
	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, 0, len(inventory))
}
