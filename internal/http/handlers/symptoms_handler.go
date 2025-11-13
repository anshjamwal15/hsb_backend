package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SymptomsHandler struct {
	symptomsService *services.SymptomsService
}

func NewSymptomsHandler(symptomsService *services.SymptomsService) *SymptomsHandler {
	return &SymptomsHandler{
		symptomsService: symptomsService,
	}
}

func (h *SymptomsHandler) SubmitTracking(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var tracking entities.SymptomsTracking
	if err := c.ShouldBindJSON(&tracking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	tracking.UserID = userOID

	if err := h.symptomsService.SubmitTracking(c.Request.Context(), &tracking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tracking,
	})
}

func (h *SymptomsHandler) GetTrackingHistory(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	trackings, err := h.symptomsService.GetTrackingHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trackings,
	})
}
