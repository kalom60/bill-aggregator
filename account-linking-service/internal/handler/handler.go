package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/database"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/grpc"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Handler interface {
	LinkAccount(c *gin.Context)
	GetLinkedAccountsByUserID(c *gin.Context)
	DeleteLinkedAccounts(c *gin.Context)
}

type handler struct {
	client     database.Service
	gRPCClient *grpc.ProviderClient
}

func NewHandler(client database.Service, gRPCClient *grpc.ProviderClient) Handler {
	return &handler{
		client:     client,
		gRPCClient: gRPCClient,
	}
}

func (h *handler) LinkAccount(c *gin.Context) {
	validatedAccount, exists := c.Get("validatedAccount")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Validation failed"})
		return
	}

	req := validatedAccount.(models.Account)

	exists, err := h.gRPCClient.IsProviderExist(req.ProviderID)
	if err != nil {
		fmt.Println("error 1", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to validate provider"})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Provider does not exist"})
		return
	}

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
