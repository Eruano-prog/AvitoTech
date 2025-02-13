package postgres

import (
	"AvitoTech/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestInsertOperation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operation := entity.Operation{
		FromUser: "user1",
		ToUser:   "user2",
		Amount:   100,
	}

	insertedOperation, err := repo.InsertOperation(operation)
	assert.NoError(t, err)
	assert.NotNil(t, insertedOperation)
	assert.Equal(t, "user1", insertedOperation.FromUser)
	assert.Equal(t, "user2", insertedOperation.ToUser)
	assert.Equal(t, 100, insertedOperation.Amount)
	defer repo.DeleteOperation(insertedOperation.Id)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM history WHERE sender_name = $1 AND receiver_name = $2 AND amount = $3",
		operation.FromUser, operation.ToUser, operation.Amount).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetSentByUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operation := entity.Operation{
		FromUser: "user1",
		ToUser:   "user2",
		Amount:   100,
	}
	o, err := repo.InsertOperation(operation)
	assert.NoError(t, err)
	defer repo.DeleteOperation(o.Id)

	operations, err := repo.GetSentByUser("user1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(operations))
	assert.Equal(t, "user1", operations[0].FromUser)
	assert.Equal(t, "user2", operations[0].ToUser)
	assert.Equal(t, 100, operations[0].Amount)
}

func TestGetReceivedByUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operation := entity.Operation{
		FromUser: "user1",
		ToUser:   "user2",
		Amount:   100,
	}
	o, err := repo.InsertOperation(operation)
	assert.NoError(t, err)
	defer repo.DeleteOperation(o.Id)

	operations, err := repo.GetReceivedByUser("user2")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(operations))
	assert.Equal(t, "user1", operations[0].FromUser)
	assert.Equal(t, "user2", operations[0].ToUser)
	assert.Equal(t, 100, operations[0].Amount)
}

func TestGetSentByUserEmpty(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operations, err := repo.GetSentByUser("nonexistent_user")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(operations))
}

func TestGetReceivedByUserEmpty(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operations, err := repo.GetReceivedByUser("nonexistent_user")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(operations))
}

func TestDeleteOperation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	operation := entity.Operation{
		FromUser: "user1",
		ToUser:   "user2",
		Amount:   100,
	}
	insertedOperation, err := repo.InsertOperation(operation)
	assert.NoError(t, err)

	err = repo.DeleteOperation(insertedOperation.Id)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM history WHERE id = $1", insertedOperation.Id).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDeleteNonExistentOperation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewHistoryRepository(logger, db)

	err := repo.DeleteOperation(99999)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM history WHERE id = $1", 99999).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
