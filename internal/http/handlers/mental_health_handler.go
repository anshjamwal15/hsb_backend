package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MentalHealthHandler struct {
	mentalHealthService *services.MentalHealthService
}

func NewMentalHealthHandler(mentalHealthService *services.MentalHealthService) *MentalHealthHandler {
	return &MentalHealthHandler{
		mentalHealthService: mentalHealthService,
	}
}

func (h *MentalHealthHandler) GetTests(c *gin.Context) {
	tests := h.mentalHealthService.GetTests()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tests,
	})
}

func (h *MentalHealthHandler) GetTestByName(c *gin.Context) {
	testName := c.Param("testName")

	test := h.mentalHealthService.GetTestByName(testName)
	if test == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Test not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    test,
	})
}

func (h *MentalHealthHandler) SubmitTestResults(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var result entities.TestResult
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result.UserID = userOID

	if err := h.mentalHealthService.SubmitTestResults(c.Request.Context(), &result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

func (h *MentalHealthHandler) GetTestResults(c *gin.Context) {
	userID := c.GetString("userID")
	testName := c.Query("testName")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	results, err := h.mentalHealthService.GetTestResults(c.Request.Context(), userOID.Hex(), testName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}
