package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeriodHandler struct {
	periodService *services.PeriodService
}

func NewPeriodHandler(periodService *services.PeriodService) *PeriodHandler {
	return &PeriodHandler{
		periodService: periodService,
	}
}

func (h *PeriodHandler) GetPeriodCycle(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	cycles, err := h.periodService.GetPeriodCycles(c.Request.Context(), userOID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cycles,
	})
}

func (h *PeriodHandler) AddPeriodCycle(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var cycle entities.PeriodCycle
	if err := c.ShouldBindJSON(&cycle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	cycle.UserID = userOID

	if err := h.periodService.AddPeriodCycle(c.Request.Context(), &cycle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    cycle,
	})
}

func (h *PeriodHandler) ResetPeriodTracker(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	if err := h.periodService.ResetPeriodTracker(c.Request.Context(), userOID.Hex()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Period tracker reset successfully",
	})
}
