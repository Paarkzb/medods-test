package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"-"`
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
}

type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type RefreshToken struct {
	Token   string    `json:"token"`
	ExpTime time.Time `json:"exp_time"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
