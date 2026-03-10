package user

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/chess-online/models"
	"github.com/ximofam/chess-online/services/auth"
	"github.com/ximofam/chess-online/services/auth/password"
	"github.com/ximofam/chess-online/services/auth/token"
	"github.com/ximofam/chess-online/services/utils"
)

type Handler struct {
	db Database
}

func NewHandler(db Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetMe(c *gin.Context) {
	user := auth.GetUser(c)
	c.JSON(200, user)
}

func (h *Handler) Login(c *gin.Context) {
	req, ok := utils.BindJSON[LoginRequest](c)
	if !ok {
		return
	}

	ctx := c.Request.Context()
	user, err := h.db.GetUserByUsername(ctx, req.Username)
	if err != nil || user == nil {
		c.JSON(401, "incorrect username or password")
		return
	}

	if !password.ComparePassword([]byte(user.Password), []byte(req.Password)) {
		c.JSON(401, "incorrect username or password")
		return
	}

	tokenRes, err := h.generateToken(ctx, user.ID)
	if err != nil {
		log.Printf("generate token error: %v", err)
		c.Status(500)
	}

	c.JSON(200, tokenRes)
}

func (h *Handler) Register(c *gin.Context) {
	req, ok := utils.BindJSON[RegisterRequest](c)
	if !ok {
		return
	}

	ctx := c.Request.Context()
	user, err := h.db.GetUserByUsername(ctx, req.Username)
	if err != nil || user != nil {
		c.JSON(400, "user already exists")
		return
	}

	hash := password.HashPassword(req.Password)

	createUser := models.User{
		Username: req.Username,
		Password: string(hash),
	}
	err = h.db.CreateUser(ctx, &createUser)
	if err != nil {
		c.Status(500)
		return
	}

	c.Status(201)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	req, ok := utils.BindJSON[RefreshTokenRequest](c)
	if !ok {
		return
	}

	ctx := c.Request.Context()
	refreshToken, err := h.db.GetRefreshTokenByToken(ctx, req.Token)
	if err != nil || refreshToken.Revoked {
		c.JSON(401, "invalid token")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		c.JSON(401, "token has expireed")
		return
	}

	if err := h.db.RevokeRefreshToken(ctx, refreshToken.ID); err != nil {
		c.JSON(401, "failed to refresh token")
		return
	}

	tokenRes, err := h.generateToken(ctx, refreshToken.UserID)
	if err != nil {
		log.Printf("generate token error: %v", err)
		c.Status(500)
	}

	c.JSON(200, tokenRes)
}

func (h *Handler) generateToken(ctx context.Context, userID uint) (*TokenResponse, error) {
	accessTokenStr, err := token.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	if err := h.db.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshToken.Token,
	}, nil
}
