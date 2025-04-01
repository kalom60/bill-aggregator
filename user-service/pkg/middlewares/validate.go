package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kalom60/bill-aggregator/user-service/internal/models"
)

func SignupMiddleware(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid request body"})
		c.Abort()
		return
	}

	user.Email = strings.TrimSpace(user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	validate := validator.New()

	_ = validate.RegisterValidation("matches", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "" || strings.ContainsAny(fl.Field().String(), "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz ")
	})

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		c.Abort()
		return
	}

	c.Set("validatedUser", user)

	c.Next()
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func LoginMiddleware(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid request body"})
		c.Abort()
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		c.Abort()
		return
	}

	c.Set("validatedRequest", req)

	c.Next()
}
