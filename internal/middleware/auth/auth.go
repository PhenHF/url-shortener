package middlewareAuth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	config "github.com/PhenHF/url-shortener/internal/config"
	middlewareLogger "github.com/PhenHF/url-shortener/internal/middleware/logger"
	storage "github.com/PhenHF/url-shortener/internal/storage"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var TOKEN_EXP = time.Hour * 3

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("JWT")
		fmt.Println(token)
		if err != nil {
			userID, err := storage.ReadyStorage.CreateUser(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				middlewareLogger.Log.Error("ERROR", zap.String("msg", err.Error()))
				return
			}

			token, err := buildJWTstring(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				middlewareLogger.Log.Error("ERROR", zap.String("msg", err.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			cookie := &http.Cookie{Name: "JWT", Value: token}
			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		userID := getUserID(token.Value)
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func buildJWTstring(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenStr, err := token.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		return "", nil
	}

	return tokenStr, nil
}

func getUserID(tokenStr string) uint {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.SECRET_KEY), nil
	})
	if err != nil {
		return 0
	}

	if !token.Valid {
		return 0
	}

	return claims.UserID
}
