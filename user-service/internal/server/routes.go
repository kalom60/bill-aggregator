package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/user-service/pkg/middlewares"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/register", middlewares.SignupMiddleware, s.handler.Register)
	r.POST("/login", middlewares.LoginMiddleware, s.handler.Login)

	return r
}
