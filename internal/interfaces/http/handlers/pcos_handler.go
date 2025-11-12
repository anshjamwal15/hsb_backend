package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PCOSHandler struct {
	pcosService *services.PCOSService
}

func NewPCOSHandler(pcosService *services.PCOSService) *PCOSHandler {
	return &PCOSHandler{
		pcosService: pcosService,
	}
}

func (h *PCOSHandler) GetQuestions(c *gin.Context) {
	questions := h.pcosService.GetQuestions()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    questions,
	})
}

func (h *PCOSHandler) SubmitAssessment(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var assessment entities.PCOSAssessment
	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	assessment.UserID = userOID

	if err := h.pcosService.SubmitAssessment(c.Request.Context(), &assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    assessment,
	})
}

func (h *PCOSHandler) GetHistory(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	assessments, err := h.pcosService.GetHistory(c.Request.Context(), userOID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assessments,
	})
}

func (h *PCOSHandler) GetLatestAssessment(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	assessment, err := h.pcosService.GetLatestAssessment(c.Request.Context(), userOID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assessment,
	})
}
