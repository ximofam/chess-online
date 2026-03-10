package database

import (
	"context"
	"errors"

	"github.com/ximofam/chess-online/models"
	"gorm.io/gorm"
)

func (g *GormDatabase) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := g.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (g *GormDatabase) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := g.DB.WithContext(ctx).Where("username = ?", username).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (g *GormDatabase) CreateUser(ctx context.Context, user *models.User) error {
	return g.DB.WithContext(ctx).Create(user).Error
}

func (g *GormDatabase) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	return g.DB.WithContext(ctx).Create(refreshToken).Error
}

func (g *GormDatabase) GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := g.DB.WithContext(ctx).Where("token = ?", token).Take(&refreshToken).Error
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (g *GormDatabase) RevokeRefreshToken(ctx context.Context, id uint) error {
	return g.DB.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ? AND revoked = ?", id, false).
		Update("revoked", true).Error
}

func (g *GormDatabase) GetUsersByIDs(ctx context.Context, ids []uint) ([]models.User, error) {
	var users []models.User
	err := g.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
