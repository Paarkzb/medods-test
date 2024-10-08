package repository

import (
	"context"
	"medodstest/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := "INSERT INTO public.user (name, email, username, password) VALUES ($1, $2, $3, $4) RETURNING id;"

	err := r.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Username, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (model.User, error) {
	var user model.User

	query := `
				SELECT u.id
				FROM public.user as u 
				WHERE u.username = $1 AND u.password = $2
			`

	err := r.db.QueryRow(context.Background(), query, username, password).Scan(&user.Id)

	return user, err
}

func (r *AuthPostgres) GetUserById(id uuid.UUID) (model.UserResponse, error) {
	var user model.UserResponse

	query := `
				SELECT u.id, u.name, u.email, u.username
				FROM public.user as u 
				WHERE u.id = $1
	`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Username)

	return user, err
}

func (r *AuthPostgres) SetRefreshToken(userId uuid.UUID, refreshToken model.RefreshToken) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
		UPDATE public.user SET refresh_token=$1, refresh_expiration_time=$2 WHERE id=$3
	`

	_, err = tx.Exec(context.Background(), query, refreshToken.Token, refreshToken.ExpTime, userId)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (r *AuthPostgres) GetRefreshToken(userId uuid.UUID) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	query := `
				SELECT u.refresh_token, u.refresh_expiration_time
				FROM public.user as u 
				WHERE u.id = $1
	`

	err := r.db.QueryRow(context.Background(), query, userId).Scan(&refreshToken.Token, &refreshToken.ExpTime)

	return refreshToken, err
}

func (r *AuthPostgres) SetIpAddress(userId uuid.UUID, ipAddress string) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
		UPDATE public.user SET ip_address=$1 WHERE id=$2
	`

	_, err = tx.Exec(context.Background(), query, ipAddress, userId)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (r *AuthPostgres) GetIpAddress(userId uuid.UUID) (string, error) {
	var ipAddress string

	query := `
				SELECT u.ip_address
				FROM public.user as u 
				WHERE u.id = $1
	`

	err := r.db.QueryRow(context.Background(), query, userId).Scan(&ipAddress)

	return ipAddress, err
}
