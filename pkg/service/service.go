package service

import (
	"medodstest/internal/model"
	"medodstest/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(userId int) (model.UserResponse, error)
	GenerateAccessToken(username, password string) (string, error)
	GenerateRefreshToken(userId int) (string, error)
	ParseToken(accessToken string) (int, error)
	RefreshTokens(userId int, refreshToken string) (model.Tokens, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
