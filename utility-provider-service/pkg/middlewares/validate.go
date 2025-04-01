package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/models"
)

func CreateProviderMiddleware(c *gin.Context) {
	var provider models.Provider

	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid request body"})
		c.Abort()
		return
	}

	provider.Name = strings.TrimSpace(provider.Name)
	provider.API_URL = strings.TrimSpace(provider.API_URL)
	provider.Authentication_Type = strings.TrimSpace(provider.Authentication_Type)
	provider.API_key = strings.TrimSpace(provider.API_key)

	validate := validator.New()

	if err := validate.Struct(provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		c.Abort()
		return
	}

	c.Set("validatedProvider", provider)

	c.Next()
}
