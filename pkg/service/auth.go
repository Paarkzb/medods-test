package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"medodstest/internal/model"
	"medodstest/pkg/repository"
	"time"

	"math/rand"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt            = "kjasdhflkqwurh"
	signingKey      = "askdjfsa;ldfkjdsal;128"
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 72 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"sub"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GetUser(userId int) (model.UserResponse, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *AuthService) GenerateAccessToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	return generateAccessTokenById(user.Id)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неправильный метод подписи")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("тип токена не совпадает с типом tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) GenerateRefreshToken(userId int) (string, error) {

	b := make([]byte, 32)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	expTime := time.Now().Add(refreshTokenTTL)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	hashToken, err := bcrypt.GenerateFromPassword(b, 10)
	if err != nil {
		return "", err
	}

	err = s.repo.SetRefreshToken(userId, model.RefreshToken{Token: string(hashToken), ExpTime: expTime})
	if err != nil {
		return "", err
	}

	b64Token := base64.StdEncoding.EncodeToString(hashToken)

	return b64Token, nil
}

func generateAccessTokenById(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		&tokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			userId,
		})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) RefreshTokens(userId int, refreshToken string) (model.Tokens, error) {
	var tokens model.Tokens

	dbRefreshToken, err := s.repo.GetRefreshToken(userId)
	if err != nil {
		return tokens, err
	}

	bRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return tokens, err
	}

	if dbRefreshToken.Token != string(bRefreshToken) {
		return tokens, errors.New("Рефреш токен неверный")
	}

	if time.Now().After(dbRefreshToken.ExpTime) {
		return tokens, errors.New("Время действия рефреш токена истекло")
	}

	accessToken, err := generateAccessTokenById(userId)
	if err != nil {
		return tokens, err
	}

	newRefreshToken, err := s.GenerateRefreshToken(userId)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = newRefreshToken

	return tokens, nil
}
