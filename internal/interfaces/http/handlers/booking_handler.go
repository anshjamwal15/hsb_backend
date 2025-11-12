package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		DoctorID    string `json:"doctorId" binding:"required"`
		SessionType string `json:"sessionType" binding:"required"`
		Date        string `json:"date" binding:"required"`
		TimeSlot    string `json:"timeSlot" binding:"required"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid date format"})
		return
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), userID, req.DoctorID, req.SessionType, date, req.TimeSlot, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"bookingId":       booking.ID.Hex(),
			"razorpayOrderId": booking.RazorpayOrderID,
			"amount":          booking.Amount,
			"status":          booking.Status,
		},
	})
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	bookings, total, err := h.bookingService.GetUserBookings(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (h *BookingHandler) GetActiveBookings(c *gin.Context) {
	userID := c.GetString("userID")

	bookings, err := h.bookingService.GetActiveBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}

func (h *BookingHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		BookingID         string `json:"bookingId"`
		RazorpayOrderID   string `json:"razorpayOrderId" binding:"required"`
		RazorpayPaymentID string `json:"razorpayPaymentId" binding:"required"`
		RazorpaySignature string `json:"razorpaySignature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.bookingService.VerifyPayment(c.Request.Context(), req.BookingID, req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment verified successfully",
	})
}
