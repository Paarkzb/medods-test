package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		headerParts := strings.Split(header["Authorization"][0], " ")
		if len(headerParts) == 0 {
			newErrorResponse(w, http.StatusUnauthorized, "пустой header авторизации")
			return
		}

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			newErrorResponse(w, http.StatusUnauthorized, "неправильный header авторизации")
			return
		}

		if len(headerParts[1]) == 0 {
			newErrorResponse(w, http.StatusUnauthorized, "пустой токен")
			return
		}

		userId, err := h.service.Authorization.ParseToken(headerParts[1])
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserId(r *http.Request) (int, error) {
	id := r.Context().Value(userCtx)

	idint, ok := id.(int)
	if !ok {
		return 0, errors.New("неправильный тип у userId")
	}

	return idint, nil
}
