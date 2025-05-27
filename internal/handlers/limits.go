package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dnakolan/rate-limiter-service/internal/models"
	"github.com/dnakolan/rate-limiter-service/internal/services"
)

type LimitsHandler struct {
	service services.LimitsService
}

func NewLimitsHandler(service services.LimitsService) *LimitsHandler {
	return &LimitsHandler{
		service: service,
	}
}

func (h *LimitsHandler) CreateRateLimitHandler(c *gin.Context) {
	var req models.CreateRateLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRateLimit(c.Request.Context(), &req.RateLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, req.RateLimit)
}
