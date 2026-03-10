package user

import (
	"context"

	"github.com/ximofam/chess-online/models"
)

type Database interface {
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uint) error
}
