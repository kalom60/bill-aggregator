package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/database"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Handler interface {
	LinkAccount(c *gin.Context)
	GetLinkedAccountsByUserID(c *gin.Context)
	DeleteLinkedAccounts(c *gin.Context)
}

type handler struct {
	client database.Service
}

func NewHandler(client database.Service) Handler {
	return &handler{
		client: client,
	}
}

func (h *handler) LinkAccount(c *gin.Context) {
	validatedAccount, exists := c.Get("validatedAccount")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Validation failed"})
		return
	}

	req := validatedAccount.(models.Account)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.EncryptedCredential), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to hash credential"})
		return
	}

	err = h.client.InsertOne(req.UserID, req.ProviderID, req.AccountIdentifier, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to link account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Account linked successfully",
	})
}

func (h *handler) GetLinkedAccountsByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	accounts, err := h.client.FetchLinkedAccountsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Failed fetching accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Accounts fetched successfully",
		"data":    accounts,
	})
}

func (h *handler) DeleteLinkedAccounts(c *gin.Context) {
	accountID := c.Param("account_id")

	err := h.client.DeleteLinkedAccountByID(accountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Account deleted successfully",
	})
}
