package service

import (
	"medodstest/internal/model"
	"medodstest/pkg/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user model.User) (uuid.UUID, error)
	GetUser(userId uuid.UUID) (model.UserResponse, error)
	GenerateAccessToken(username, password string) (string, error)
	GenerateRefreshToken(userId uuid.UUID, ipAddress string) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
	RefreshTokens(userId uuid.UUID, refreshToken, ipAddress string) (model.Tokens, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
