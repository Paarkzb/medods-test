package model

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
