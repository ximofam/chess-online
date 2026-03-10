package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ximofam/chess-online/config"
	"github.com/ximofam/chess-online/models"
)

var (
	once            sync.Once
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
)

func Init(cfg config.JWTConfig) {
	once.Do(func() {
		secretKey = []byte(cfg.SecretKey)
		accessTokenTTL = cfg.AccessTokenTTL
		refreshTokenTTL = cfg.RefreshTokenTTL
	})
}

type AccessClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userID uint) (*models.RefreshToken, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	token := hex.EncodeToString(b)

	return &models.RefreshToken{
		Token:     token,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
		UserID:    userID,
	}, nil
}

func ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	claims := AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}
