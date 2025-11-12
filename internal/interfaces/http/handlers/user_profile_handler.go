package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	userService *services.UserService
}

func NewUserProfileHandler(userService *services.UserService) *UserProfileHandler {
	return &UserProfileHandler{
		userService: userService,
	}
}

// GetProfile retrieves the authenticated user's profile
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "user not found",
		})
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

// UpdateProfile updates the authenticated user's profile
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	var req struct {
		Name         string `json:"name"`
		PhoneNumber  string `json:"phoneNumber"`
		ProfileImage string `json:"profileImage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req.Name, req.PhoneNumber, req.ProfileImage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
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
