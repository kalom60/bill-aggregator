package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/broker/pkg/middlewares"
)

type RequestPayload struct {
	Action              string                     `json:"action"`
	UserRegister        UserRegisterPayload        `json:"user_register,omitempty"`
	UserLogin           UserLoginPayload           `json:"user_login,omitempty"`
	CreateProvider      CreateProviderPayload      `json:"create_provider,omitempty"`
	LinkAccount         LinkAccountPayload         `json:"link_account,omitempty"`
	DeleteLinkedAccount DeleteLinkedAccountPayload `json:"delete_linked_account,omitempty"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UserRegisterPayload struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateProviderPayload struct {
	Name               string `json:"name"`
	APIURL             string `json:"api_url"`
	AuthenticationType string `json:"authentication_type"`
	APIKey             string `json:"api_key"`
}

type LinkAccountPayload struct {
	UserID              string `json:"user_id,omitempty"`
	ProviderID          string `json:"provider_id"`
	AccountIdentifier   string `json:"account_identifier"`
	EncryptedCredential string `json:"encrypted_credential"`
}

type DeleteLinkedAccountPayload struct {
	AccountID string `json:"account_id"`
}

func (app *Config) Broker(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Hit the broker"})
}

func (app *Config) HandleSubmission(c *gin.Context) {
	var req RequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid request payload"})
		return
	}

	if req.Action != "user_register" && req.Action != "user_login" {
		middlewares.Authenticate(c)
		if c.IsAborted() {
			return
		}
	}

	switch req.Action {
	case "user_register":
		app.UserRegister(c, req.UserRegister)
	case "user_login":
		app.UserLogin(c, req.UserLogin)
	case "create_provider":
		app.CreateProvider(c, req.CreateProvider)
	case "get_provider":
		app.GetProvider(c)
	case "link_account":
		user_id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Authentication failed"})
			return
		}

		if userIDStr, ok := user_id.(string); ok {
			req.LinkAccount.UserID = userIDStr
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Invalid user ID format"})
			return
		}

		app.LinkAccount(c, req.LinkAccount)
	case "get_linked_accounts":
		user_id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Authentication failed"})
			return
		}

		userIDStr, ok := user_id.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Invalid user ID format"})
			return
		}

		app.GetLinkedAcounts(c, userIDStr)
	case "delete_linked_account":
		app.DeleteLinkedAccount(c, req.DeleteLinkedAccount)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Unknown action"})
	}
}

func (app *Config) UserRegister(c *gin.Context, payload UserRegisterPayload) {
	app.forwardRequest(c, "POST", "http://user-service/register", payload)
}

func (app *Config) UserLogin(c *gin.Context, payload UserLoginPayload) {
	app.forwardRequest(c, "POST", "http://user-service/login", payload)
}

func (app *Config) CreateProvider(c *gin.Context, payload CreateProviderPayload) {
	app.forwardRequest(c, "POST", "http://utility-provider-service/provider", payload)
}

func (app *Config) GetProvider(c *gin.Context) {
	app.forwardRequest(c, "GET", "http://utility-provider-service/provider", nil)
}

func (app *Config) LinkAccount(c *gin.Context, payload LinkAccountPayload) {
	app.forwardRequest(c, "POST", "http://account-linking-service/accounts/link", payload)
}

func (app *Config) GetLinkedAcounts(c *gin.Context, userID string) {
	url := fmt.Sprintf("http://account-linking-service/accounts/%s", userID)
	app.forwardRequest(c, "GET", url, nil)
}

func (app *Config) DeleteLinkedAccount(c *gin.Context, payload DeleteLinkedAccountPayload) {
	url := fmt.Sprintf("http://account-linking-service/accounts/%s", payload.AccountID)
	app.forwardRequest(c, "DELETE", url, nil)
}

func (app *Config) forwardRequest(c *gin.Context, method, url string, payload any) {
	var reqBody *bytes.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to encode request"})
			return
		}
		reqBody = bytes.NewReader(jsonData)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	var jsonResponse JsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to decode response"})
		return
	}

	if jsonResponse.Error {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": jsonResponse.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Success", "data": jsonResponse.Data})
}
