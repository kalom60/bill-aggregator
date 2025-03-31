package server

import (
	"net/http"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/register", middlewares.SignupMiddleware, s.handler.Register)
	r.POST("/login", middlewares.LoginMiddleware, s.handler.Login)

	return r
}
