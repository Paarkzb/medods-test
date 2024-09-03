package repository

import (
	"medodstest/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(user model.User) (uuid.UUID, error)
	GetUser(username, password string) (model.User, error)
	GetUserById(id uuid.UUID) (model.UserResponse, error)
	SetRefreshToken(userId uuid.UUID, refreshToken model.RefreshToken) error
	GetRefreshToken(userId uuid.UUID) (model.RefreshToken, error)
	SetIpAddress(userId uuid.UUID, ipAddress string) error
	GetIpAddress(userId uuid.UUID) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
