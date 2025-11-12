package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PregnancyHandler struct {
	pregnancyService *services.PregnancyService
}

func NewPregnancyHandler(pregnancyService *services.PregnancyService) *PregnancyHandler {
	return &PregnancyHandler{
		pregnancyService: pregnancyService,
	}
}

func (h *PregnancyHandler) GetPregnancyData(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	tracker, err := h.pregnancyService.GetPregnancyData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tracker,
	})
}

func (h *PregnancyHandler) AddPregnancyEntry(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var tracker entities.PregnancyTracker
	if err := c.ShouldBindJSON(&tracker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	tracker.UserID = userOID

	if err := h.pregnancyService.AddPregnancyEntry(c.Request.Context(), &tracker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tracker,
	})
}
