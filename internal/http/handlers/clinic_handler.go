package handlers

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	clinicService *services.ClinicService
}

func NewClinicHandler(clinicService *services.ClinicService) *ClinicHandler {
	return &ClinicHandler{
		clinicService: clinicService,
	}
}

func (h *ClinicHandler) GetClinics(c *gin.Context) {
	clinics, err := h.clinicService.GetClinics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clinics,
	})
}

func (h *ClinicHandler) CreateBooking(c *gin.Context) {
	userID := c.GetString("userID")

	var booking entities.ClinicBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	booking.UserID = userID

	if err := h.clinicService.CreateBooking(c.Request.Context(), &booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    booking,
	})
}

func (h *ClinicHandler) GetMyBookings(c *gin.Context) {
	userID := c.GetString("userID")

	bookings, err := h.clinicService.GetMyBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}

func (h *ClinicHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		BookingID string `json:"bookingId" binding:"required"`
		PaymentID string `json:"paymentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.clinicService.VerifyPayment(c.Request.Context(), req.BookingID, req.PaymentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment verified successfully",
	})
}
