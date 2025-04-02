package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/account-linking-service/pkg/middlewares"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/accounts/link", middlewares.LinkAccountMiddleware, s.handler.LinkAccount)
	r.GET("/accounts/:user_id", s.handler.GetLinkedAccountsByUserID)
	r.DELETE("/accounts/:account_id", s.handler.DeleteLinkedAccounts)

	return r
}
