package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type JournalHandler struct {
	journalService *services.JournalService
}

func NewJournalHandler(journalService *services.JournalService) *JournalHandler {
	return &JournalHandler{
		journalService: journalService,
	}
}

type JournalRequest struct {
	UserID   string `json:"userId" binding:"required"`
	Title    string `json:"title"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category" binding:"required"`
}

func (h *JournalHandler) CreateJournal(c *gin.Context) {
	var req JournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	journal, err := h.journalService.CreateJournal(c.Request.Context(), req.UserID, req.Title, req.Content, req.Category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    journal,
	})
}

func (h *JournalHandler) GetJournals(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")
	dateStr := c.Query("date")

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = &parsedDate
		}
	}

	journals, total, err := h.journalService.GetJournals(c.Request.Context(), userID, page, limit, search, category, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    journals,
		"total":   total,
		"page":    page,
		"pages":   totalPages,
	})
}

func (h *JournalHandler) GetJournalByID(c *gin.Context) {
	journalID := c.Param("journalId")

	journal, err := h.journalService.GetJournalByID(c.Request.Context(), journalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Journal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    journal,
	})
}

func (h *JournalHandler) UpdateJournal(c *gin.Context) {
	journalID := c.Param("journalId")

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	journal, err := h.journalService.UpdateJournal(c.Request.Context(), journalID, req.Title, req.Content, req.Category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    journal,
	})
}

func (h *JournalHandler) DeleteJournal(c *gin.Context) {
	journalID := c.Param("journalId")

	if err := h.journalService.DeleteJournal(c.Request.Context(), journalID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Journal deleted successfully",
	})
}

func (h *JournalHandler) GetJournalsByUserID(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	journals, total, err := h.journalService.GetJournals(c.Request.Context(), userID, page, limit, "", "", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    journals,
		"total":   total,
		"page":    page,
	})
}
