package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/broker/pkg/middlewares"
)

type RequestPayload struct {
	Action         string                `json:"action"`
	UserRegister   UserRegisterPayload   `json:"user_register,omitempty"`
	UserLogin      UserLoginPayload      `json:"user_login,omitempty"`
	CreateProvider CreateProviderPayload `json:"create_provider,omitempty"`
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
