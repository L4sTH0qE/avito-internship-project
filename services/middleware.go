package services

import (
	u "awesomeProject/utils"
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

// JwtAuthentication - Мидлваря для проверки jwt-токена перед всеми запросами, в которых это требуется.
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/auth"}
		requestPath := r.URL.Path

		// Проверка, не требует ли запрос аутентификации.
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Получение заголовка с токеном.
		tokenHeader := r.Header.Get("Authorization")

		// Если заголовок отсутствует, то возвращаем 400.
		if tokenHeader == "" {
			u.RespondWithError(w, http.StatusBadRequest, "Missing auth token")
			return
		}

		// Получение JWT-Токена из заголовка.
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			u.RespondWithError(w, http.StatusBadRequest, "Invalid auth token")
			return
		}

		tokenPart := splitted[1]
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		// Если неправильный токен, то возвращаем 401.
		if err != nil {
			u.RespondWithError(w, http.StatusUnauthorized, "Malformed auth token")
			return
		}

		// Если недействительный токен, то возвращаем 401.
		if !token.Valid {
			u.RespondWithError(w, http.StatusUnauthorized, "Token is not valid.")
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.Username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
