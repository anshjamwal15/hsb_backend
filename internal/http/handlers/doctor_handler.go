package handlers

import (
	"net/http"
	"strconv"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorService *services.DoctorService
}

func NewDoctorHandler(doctorService *services.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
	}
}

func (h *DoctorHandler) GetDoctors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	doctors, total, err := h.doctorService.GetDoctors(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"data":              doctors,
		"current_page":      page,
		"total_pages":       totalPages,
		"total_doctors":     total,
		"has_next_page":     page < totalPages,
		"has_previous_page": page > 1,
	})
}

func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	doctorID := c.Param("doctorId")

	doctor, err := h.doctorService.GetDoctorByID(c.Request.Context(), doctorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Doctor not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    doctor,
	})
}
