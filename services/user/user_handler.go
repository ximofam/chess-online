package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/chess-online/services/utils"
)

func (h *Handler) GetUser(c *gin.Context) {
	userIDStr, ok := utils.GetPathValueFromRequest(c, "id")
	if !ok {
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, "Invalid user id")
		return
	}

	user, err := h.db.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, user)
}
