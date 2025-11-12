package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FSFIHandler struct {
	fsfiService *services.FSFIService
}

func NewFSFIHandler(fsfiService *services.FSFIService) *FSFIHandler {
	return &FSFIHandler{
		fsfiService: fsfiService,
	}
}

func (h *FSFIHandler) GetTest(c *gin.Context) {
	test := h.fsfiService.GetTest()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    test,
	})
}

func (h *FSFIHandler) SubmitTest(c *gin.Context) {
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

	if err := h.fsfiService.SubmitTest(c.Request.Context(), &result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

func (h *FSFIHandler) GetMyResults(c *gin.Context) {
	userID := c.GetString("userID")

	// Convert string userID to ObjectID
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	results, err := h.fsfiService.GetMyResults(c.Request.Context(), userOID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}
