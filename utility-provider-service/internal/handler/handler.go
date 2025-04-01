package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/database"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/models"
)

type Handler interface {
	CreateProvider(c *gin.Context)
	GetProviders(c *gin.Context)
}

type handler struct {
	client database.Service
}

func NewHandler(client database.Service) Handler {
	return &handler{
		client: client,
	}
}

func (h *handler) CreateProvider(c *gin.Context) {
	validatedProvider, exists := c.Get("validatedProvider")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Validation failed"})
		return
	}

	req := validatedProvider.(models.Provider)

	err := h.client.InsertOne(req.Name, req.API_URL, req.Authentication_Type, req.API_key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to create provider"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Provider registered successfully",
	})
}

func (h *handler) GetProviders(c *gin.Context) {
	providers, err := h.client.FetchProviders()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Failed fetching providers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "provider fetched successfully",
		"data":    providers,
	})
}
