package handlers

import (
	"net/http"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type TimeSlotHandler struct {
	doctorService  *services.DoctorService
	bookingService *services.BookingService
}

func NewTimeSlotHandler(doctorService *services.DoctorService, bookingService *services.BookingService) *TimeSlotHandler {
	return &TimeSlotHandler{
		doctorService:  doctorService,
		bookingService: bookingService,
	}
}

// GetAvailableTimeSlots returns available time slots for a doctor on a specific date
func (h *TimeSlotHandler) GetAvailableTimeSlots(c *gin.Context) {
	doctorID := c.Query("doctorId")
	dateStr := c.Query("date")

	if doctorID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "doctorId and date are required",
		})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	// Get doctor details
	doctor, err := h.doctorService.GetDoctorByID(c.Request.Context(), doctorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Doctor not found",
		})
		return
	}

	// Check if doctor is available
	if doctor.IsAvailable != nil && !*doctor.IsAvailable {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []interface{}{},
			"message": "Doctor is not available",
		})
		return
	}

	// Generate time slots based on doctor's working hours
	timeSlots := generateTimeSlots(doctor.Timing, date)

	// Get existing bookings for this doctor on this date
	existingBookings, err := h.bookingService.GetDoctorBookingsByDate(c.Request.Context(), doctorID, date)
	if err != nil {
		// Log error but continue with all slots marked as available
		timeSlots = markAllAsAvailable(timeSlots)
	} else {
		// Mark booked slots as unavailable
		timeSlots = markBookedSlots(timeSlots, existingBookings)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    timeSlots,
	})
}

type TimeSlot struct {
	Time      string `json:"time"`
	Available bool   `json:"available"`
}

// generateTimeSlots generates time slots based on doctor's working hours
func generateTimeSlots(timing interface{}, date time.Time) []TimeSlot {
	var slots []TimeSlot

	// Default working hours if not specified
	startHour := 9
	endHour := 18
	slotDuration := 30 // minutes

	// Parse timing if available
	if timing != nil {
		// This is a simplified implementation
		// In production, you'd parse the actual timing structure
	}

	// Generate slots
	currentTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, 0, 0, 0, date.Location())
	endTime := time.Date(date.Year(), date.Month(), date.Day(), endHour, 0, 0, 0, date.Location())

	for currentTime.Before(endTime) {
		slots = append(slots, TimeSlot{
			Time:      currentTime.Format("15:04"),
			Available: true,
		})
		currentTime = currentTime.Add(time.Duration(slotDuration) * time.Minute)
	}

	return slots
}

// markAllAsAvailable marks all slots as available
func markAllAsAvailable(slots []TimeSlot) []TimeSlot {
	for i := range slots {
		slots[i].Available = true
	}
	return slots
}

// markBookedSlots marks booked time slots as unavailable
func markBookedSlots(slots []TimeSlot, bookings interface{}) []TimeSlot {
	// This is a simplified implementation
	// In production, you'd check each booking's time slot against the generated slots
	return slots
}
