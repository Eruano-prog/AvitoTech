package postgres

import (
	"AvitoTech/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestInsertUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewUserRepository(logger, db)

	user := &entity.User{
		Username: "testuser",
		Password: "testpass",
		Balance:  100,
	}

	insertedUser, err := repo.InsertUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, insertedUser)
	assert.Equal(t, user.Username, insertedUser.Username)
	assert.Equal(t, user.Password, insertedUser.Password)
	assert.Equal(t, user.Balance, insertedUser.Balance)
}

func TestFindUserByUsername(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewUserRepository(logger, db)

	user := &entity.User{
		Username: "testuser",
		Password: "testpass",
		Balance:  100,
	}

	_, err := repo.InsertUser(user)
	assert.NoError(t, err)

	foundUser, err := repo.FindUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Password, foundUser.Password)
	assert.Equal(t, user.Balance, foundUser.Balance)
}

func TestFindUserById(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewUserRepository(logger, db)

	user := &entity.User{
		Username: "testuser",
		Password: "testpass",
		Balance:  100,
	}

	insertedUser, err := repo.InsertUser(user)
	assert.NoError(t, err)

	foundUser, err := repo.FindUserById(insertedUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Password, foundUser.Password)
	assert.Equal(t, user.Balance, foundUser.Balance)
}

func TestTransferMoney(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewUserRepository(logger, db)

	user1 := &entity.User{
		Username: "user1",
		Password: "pass1",
		Balance:  200,
	}

	user2 := &entity.User{
		Username: "user2",
		Password: "pass2",
		Balance:  100,
	}

	insertedUser1, err := repo.InsertUser(user1)
	assert.NoError(t, err)

	insertedUser2, err := repo.InsertUser(user2)
	assert.NoError(t, err)

	err = repo.TransferMoney(insertedUser1.Id, insertedUser2.Id, 50)
	assert.NoError(t, err)

	updatedUser1, err := repo.FindUserById(insertedUser1.Id)
	assert.NoError(t, err)
	assert.Equal(t, 150, updatedUser1.Balance)

	updatedUser2, err := repo.FindUserById(insertedUser2.Id)
	assert.NoError(t, err)
	assert.Equal(t, 150, updatedUser2.Balance)
}

func TestWithdrawMoney(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := NewUserRepository(logger, db)

	user := &entity.User{
		Username: "testuser",
		Password: "testpass",
		Balance:  200,
	}

	insertedUser, err := repo.InsertUser(user)
	assert.NoError(t, err)

	err = repo.WithdrawMoney(insertedUser.Id, 50)
	assert.NoError(t, err)

	updatedUser, err := repo.FindUserById(insertedUser.Id)
	assert.NoError(t, err)
	assert.Equal(t, 150, updatedUser.Balance)
}
