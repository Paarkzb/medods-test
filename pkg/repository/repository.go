package repository

import (
	"medodstest/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
	GetUserById(id int) (model.UserResponse, error)
	SetRefreshToken(userId int, refreshToken model.RefreshToken) error
	GetRefreshToken(userId int) (model.RefreshToken, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
