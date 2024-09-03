package handler

import (
	"encoding/json"
	"medodstest/internal/model"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input model.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.service.Authorization.GenerateAccessToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := h.service.Authorization.ParseToken(accessToken)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := h.service.GenerateRefreshToken(userId, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.Tokens{AccessToken: accessToken, RefreshToken: refreshToken})
}

type refreshInput struct {
	UserId       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var input refreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// tokens, err := h.service.Authorization.RefreshTokens(input.UserId, input.RefreshToken, strings.Split(r.RemoteAddr, ":")[0])
	tokens, err := h.service.Authorization.RefreshTokens(input.UserId, input.RefreshToken, "192.168.155.1")
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
