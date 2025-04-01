package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/broker/pkg/utils"
)

func Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Please login."})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token format."})
		c.Abort()
		return
	}

	userID, email, err := utils.VerifyToken(token)
	if err != nil {
		errorMessage := "Invalid or expired token. Please authenticate again."
		if err.Error() == "Token is expired." {
			errorMessage = "Expired token"
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": errorMessage})
		c.Abort()
		return
	}

	c.Set("userID", userID)
	c.Set("userEmail", email)

	c.Next()
}
