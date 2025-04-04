package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/bills/:user_id/provider", s.handler.GetBillsByProvider)
	r.GET("/bills/:user_id", s.handler.GetBills)
	r.POST("/bills/refresh/:user_id", s.handler.RefreshBills)

	return r
}
