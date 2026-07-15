package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenLifetime = 8 * time.Hour

type signInRequest struct {
	Password string `json:"password"`
}

type authClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// signInHandler проверяет пароль и возвращает JWT-токен
func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			errors.New("метод запроса не поддерживается"),
		)
		return
	}

	var request signInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			fmt.Errorf("ошибка чтения JSON: %w", err),
		)
		return
	}

	password := os.Getenv("TODO_PASSWORD")

	if !passwordsEqual(request.Password, password) {
		writeError(w, errors.New("неверный пароль"))
		return
	}

	token, err := createToken(password)
	if err != nil {
		writeError(
			w,
			fmt.Errorf("не удалось создать токен: %w", err),
		)
		return
	}

	writeJSON(w, map[string]string{
		"token": token,
	})
}

// auth проверяет JWT-токен перед вызовом защищённого обработчика
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		// Если пароль не задан, сервер работает без аутентификации
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, password) {
			http.Error(
				w,
				"Authentication required",
				http.StatusUnauthorized,
			)
			return
		}

		next(w, r)
	}
}

func createToken(password string) (string, error) {
	now := time.Now()

	claims := authClaims{
		PasswordHash: passwordHash(password),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(signingKey(password))
}

func validateToken(tokenValue string, password string) bool {
	if tokenValue == "" {
		return false
	}

	claims := &authClaims{}

	token, err := jwt.ParseWithClaims(
		tokenValue,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("неподдерживаемый алгоритм подписи")
			}

			return signingKey(password), nil
		},
		jwt.WithValidMethods(
			[]string{jwt.SigningMethodHS256.Alg()},
		),
	)
	if err != nil || !token.Valid {
		return false
	}

	expectedHash := passwordHash(password)

	return subtle.ConstantTimeCompare(
		[]byte(claims.PasswordHash),
		[]byte(expectedHash),
	) == 1
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func signingKey(password string) []byte {
	sum := sha256.Sum256(
		[]byte("todo-jwt-signing-key:" + password),
	)

	return sum[:]
}

func passwordsEqual(first string, second string) bool {
	firstHash := sha256.Sum256([]byte(first))
	secondHash := sha256.Sum256([]byte(second))

	return subtle.ConstantTimeCompare(
		firstHash[:],
		secondHash[:],
	) == 1
}