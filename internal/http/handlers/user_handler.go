package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":           user.ID.Hex(),
			"name":         user.Name,
			"email":        user.Email,
			"phoneNumber":  user.PhoneNumber,
			"profileImage": user.ProfileImage,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		Name         string `json:"name"`
		PhoneNumber  string `json:"phoneNumber"`
		ProfileImage string `json:"profileImage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req.Name, req.PhoneNumber, req.ProfileImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":           user.ID.Hex(),
			"name":         user.Name,
			"email":        user.Email,
			"phoneNumber":  user.PhoneNumber,
			"profileImage": user.ProfileImage,
		},
	})
}
