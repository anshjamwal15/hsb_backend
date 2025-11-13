package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

type DiagnosticHandler struct {
	diagnosticService *services.DiagnosticService
}

func NewDiagnosticHandler(diagnosticService *services.DiagnosticService) *DiagnosticHandler {
	return &DiagnosticHandler{
		diagnosticService: diagnosticService,
	}
}

func (h *DiagnosticHandler) GetDiagnostics(c *gin.Context) {
	diagnostics, err := h.diagnosticService.GetDiagnostics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    diagnostics,
	})
}

func (h *DiagnosticHandler) GetDiagnosticsUsers(c *gin.Context) {
	userID := c.GetString("userID")

	bookings, err := h.diagnosticService.GetDiagnosticsUsers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}

func (h *DiagnosticHandler) CreateBooking(c *gin.Context) {
	userID := c.GetString("userID")

	var booking entities.DiagnosticBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	booking.UserID = userID

	if err := h.diagnosticService.CreateBooking(c.Request.Context(), &booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    booking,
	})
}

func (h *DiagnosticHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		BookingID string `json:"bookingId" binding:"required"`
		PaymentID string `json:"paymentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.diagnosticService.VerifyPayment(c.Request.Context(), req.BookingID, req.PaymentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment verified successfully",
	})
}
