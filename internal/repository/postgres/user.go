package postgres

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type UserRepository struct {
	l  *zap.Logger
	db *sql.DB
}

func (u UserRepository) InsertUser(user *entity.User) (*entity.User, error) {
	q, err := u.db.Prepare(`
	INSERT INTO users (username, password, balance)
	VALUES ($1, $2, $3)
	RETURNING user_id, username, password, balance
	`)
	if err != nil {
		u.l.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			u.l.Error("Failed to close DB", zap.Error(err))
		}
	}(q)

	res := q.QueryRow(user.Username, user.Password, user.Balance)

	if res.Err() != nil {
		u.l.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.ID, &resUser.Username, &resUser.Password, &resUser.Balance)
	if err != nil {
		u.l.Error("Failed to scan inserted user", zap.Error(err))
		return nil, err
	}

	return &resUser, nil
}

func (u UserRepository) FindUserByUsername(username string) (*entity.User, error) {
	q, err := u.db.Prepare(`
	SELECT user_id, username, password, balance
	FROM users
	WHERE username = $1
`)
	if err != nil {
		u.l.Error("Failed to prepare query to find user by username", zap.String("username", username))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			u.l.Error("Failed to close DB", zap.Error(err))
		}
	}(q)

	res := q.QueryRow(username)
	if res.Err() != nil {
		u.l.Error("Failed to find user by username", zap.String("username", username))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.ID, &resUser.Username, &resUser.Password, &resUser.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrorUserNotFound
		}
		u.l.Error("Failed to scan found user by username", zap.String("username", username))
		return nil, err
	}

	return &resUser, nil
}

func (u UserRepository) FindUserByID(id int) (*entity.User, error) {
	q, err := u.db.Prepare(`
	SELECT user_id, username, password, balance
	FROM users
	WHERE user_id = $1
`)
	if err != nil {
		u.l.Error("Failed to find user by username", zap.Int("user id", id))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			u.l.Error("Failed to close DB", zap.Error(err))
		}
	}(q)

	res := q.QueryRow(id)
	if res.Err() != nil {
		u.l.Error("Failed to find user by username", zap.Int("user id", id))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.ID, &resUser.Username, &resUser.Password, &resUser.Balance)
	if err != nil {
		u.l.Error("Failed to scan found user by username", zap.Int("user id", id))
		return nil, err
	}

	return &resUser, nil
}

func (u UserRepository) TransferMoney(userFrom int, userTo int, amount int) error {
	tx, err := u.db.Begin()
	if err != nil {
		u.l.Error("Failed to begin transaction", zap.Error(err))
		return err
	}

	var balance int
	err = tx.QueryRow("SELECT balance FROM users WHERE user_id = $1 FOR UPDATE", userFrom).Scan(&balance)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		u.l.Error("Failed to check balance", zap.Error(err))
		return err
	}

	if balance < amount {
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE user_id = $2", amount, userFrom)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		u.l.Error("Failed to update balance for sender", zap.Error(err))
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE user_id = $2", amount, userTo)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		u.l.Error("Failed to update balance for receiver", zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		u.l.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (u UserRepository) WithdrawMoney(user int, amount int) error {
	tx, err := u.db.Begin()
	if err != nil {
		u.l.Error("Failed to begin transaction", zap.Error(err))
		return err
	}

	var balance int
	err = tx.QueryRow("SELECT balance FROM users WHERE user_id = $1 FOR UPDATE", user).Scan(&balance)
	if err != nil {
		u.l.Error("Failed to check balance", zap.Error(err))
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		return err
	}

	if balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE user_id = $2", amount, user)
	if err != nil {
		u.l.Error("Failed to update balance", zap.Error(err))
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		u.l.Error("Failed to commit transaction", zap.Error(err))
		err = tx.Rollback()
		if err != nil {
			u.l.Error("Failed to rollback transaction", zap.Error(err))
			return err
		}
		return err
	}

	return nil
}

func NewUserRepository(
	l *zap.Logger,
	db *sql.DB,
) repository.UserRepository {
	return &UserRepository{
		l:  l,
		db: db,
	}
}
