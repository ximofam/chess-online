package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPathValueFromRequest(c *gin.Context, key string) (string, bool) {
	v := c.Param(key)
	if v == "" {
		c.JSON(400, fmt.Sprintf("missing path value %s", key))
		return "", false
	}

	return v, true
}

func BindJSON[T any](c *gin.Context) (T, bool) {
	var res T
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(400, err.Error())
		return res, false
	}

	return res, true
}

func GenerateUUID() string {
	return uuid.New().String()
}
