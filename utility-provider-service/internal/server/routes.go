package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/utility-provider-service/pkg/middlewares"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/provider", middlewares.CreateProviderMiddleware, s.handler.CreateProvider)
	r.GET("/provider", s.handler.GetProviders)

	return r
}
