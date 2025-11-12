package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

type WeightHandler struct {
	weightService *services.WeightService
}

func NewWeightHandler(weightService *services.WeightService) *WeightHandler {
	return &WeightHandler{
		weightService: weightService,
	}
}

func (h *WeightHandler) GetData(c *gin.Context) {
	userID := c.GetString("userID")

	entries, err := h.weightService.GetData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entries,
	})
}

func (h *WeightHandler) AddEntry(c *gin.Context) {
	userID := c.GetString("userID")

	var entry entities.WeightMetabolic
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	entry.UserID = userID

	if err := h.weightService.AddEntry(c.Request.Context(), &entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    entry,
	})
}
