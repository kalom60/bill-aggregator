package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/models"
)

func LinkAccountMiddleware(c *gin.Context) {
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid request body"})
		c.Abort()
		return
	}

	account.UserID = strings.TrimSpace(account.UserID)
	account.ProviderID = strings.TrimSpace(account.ProviderID)
	account.AccountIdentifier = strings.TrimSpace(account.AccountIdentifier)
	account.EncryptedCredential = strings.TrimSpace(account.EncryptedCredential)

	validate := validator.New()

	if err := validate.Struct(account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		c.Abort()
		return
	}

	c.Set("validatedAccount", account)

	c.Next()
}
