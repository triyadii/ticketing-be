package middleware

import (
	"errors"
	"strings"
	"time"

	"ticketing-be/internal/config"
	"ticketing-be/internal/models"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uuid.UUID         `json:"user_id"`
	Username string            `json:"username"`
	Nama     string            `json:"nama"`
	Status   models.UserStatus `json:"status"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(time.Duration(cfg.JWTExpireHours) * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Nama:     user.Nama,
		Status:   user.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.AppName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau telah kadaluarsa")
	}

	return claims, nil
}

// Protected adalah middleware Fiber untuk memeriksa token JWT pada header Authorization
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Header Authorization diperlukan", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Format token harus 'Bearer <token>'", nil)
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Token autentikasi tidak valid atau telah kadaluarsa", err.Error())
		}

		if claims.Status != models.UserStatusAktif {
			return utils.SendError(c, fiber.StatusForbidden, "Akun petugas IT tidak aktif", nil)
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
