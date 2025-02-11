package repository

import (
	"AvitoTech/internal/entity"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type UserDb struct {
	l  *zap.Logger
	db *sql.DB
}

func (u UserDb) InsertUser(user *entity.User) (*entity.User, error) {
	q, err := u.db.Prepare(`
	INSERT INTO users (username, password, balance)
	VALUES ($1, $2, $3)
	RETURNING id, username, password, balance
	`)
	if err != nil {
		u.l.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}
	defer q.Close()

	res := q.QueryRow(user.Username, user.Password, user.Balance)

	if res.Err() != nil {
		u.l.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.Id, user.Username, &resUser.Password, &resUser.Balance)
	if err != nil {
		u.l.Error("Failed to scan inserted user", zap.Error(err))
		return nil, err
	}

	return &resUser, nil
}

func (u UserDb) FindUserByUsername(username string) (*entity.User, error) {
	q, err := u.db.Prepare(`
	SELECT id, username, password
	FROM users
	WHERE username = $1
`)
	if err != nil {
		u.l.Error("Failed to find user by username", zap.String("username", username))
		return nil, err
	}
	defer q.Close()

	res := q.QueryRow(username)
	if res.Err() != nil {
		u.l.Error("Failed to find user by username", zap.String("username", username))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.Id, &resUser.Username, &resUser.Password)
	if err != nil {
		u.l.Error("Failed to scan found user by username", zap.String("username", username))
		return nil, err
	}

	return &resUser, nil
}

func (u UserDb) FindUserById(id int) (*entity.User, error) {
	q, err := u.db.Prepare(`
	SELECT id, username, password
	FROM users
	WHERE id = $1
`)
	if err != nil {
		u.l.Error("Failed to find user by username", zap.Int("user id", id))
		return nil, err
	}
	defer q.Close()

	res := q.QueryRow(id)
	if res.Err() != nil {
		u.l.Error("Failed to find user by username", zap.Int("user id", id))
		return nil, err
	}

	var resUser entity.User
	err = res.Scan(&resUser.Id, &resUser.Username, &resUser.Password)
	if err != nil {
		u.l.Error("Failed to scan found user by username", zap.Int("user id", id))
		return nil, err
	}

	return &resUser, nil
}

func NewUserDatabase(
	l *zap.Logger,
	pgAddress string,
	pgUser string,
	pgPassword string,
	pgDatabase string,
) (*UserDb, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPassword, pgAddress, pgDatabase)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Fatal("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}

	return &UserDb{
		l:  l,
		db: db,
	}, nil
}
