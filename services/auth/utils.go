package auth

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/chess-online/models"
	"gorm.io/gorm"
)

var (
	once sync.Once
	gdb  *gorm.DB
)

func Init(db *gorm.DB) {
	once.Do(func() {
		gdb = db
	})
}

type contextKey string

const (
	userIDKey contextKey = "user_id"
	userKey   contextKey = "user"
)

func GetUserID(c *gin.Context) uint {
	return c.MustGet(userIDKey).(uint)
}

func GetUser(c *gin.Context) *models.User {
	u, ok := c.Get(userKey)
	if ok {
		return u.(*models.User)
	}

	userID := GetUserID(c)

	var user models.User
	if err := gdb.WithContext(c.Request.Context()).Take(&user, userID).Error; err != nil {
		panic(fmt.Sprintf("not found user with id: %d", userID))
	}

	c.Set(userKey, &user)

	return &user
}
