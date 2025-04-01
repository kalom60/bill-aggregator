package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/user-service/internal/database"
	"github.com/kalom60/bill-aggregator/user-service/internal/models"
	"github.com/kalom60/bill-aggregator/user-service/pkg/middlewares"
	"golang.org/x/crypto/bcrypt"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	client database.Service
}

func NewHandler(client database.Service) Handler {
	return &handler{
		client: client,
	}
}

func (h *handler) Register(c *gin.Context) {
	validatedUser, exists := c.Get("validatedUser")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Validation failed"})
		return
	}

	req := validatedUser.(models.User)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to hash password"})
		return
	}

	err = h.client.InsertOne(req.Email, string(hashedPassword), req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "User registered successfully",
	})
}

type login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *handler) Login(c *gin.Context) {
	validatedRequest, exists := c.Get("validatedRequest")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Validation failed"})
		return
	}

	loginRequest := validatedRequest.(middlewares.LoginRequest)

	user, err := h.client.GetUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid credentials"})
		return
	}

	token, err := middlewares.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
		},
	})
}
