package server

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/database"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/handler"
)

type Server struct {
	port int

	db      database.Service
	handler handler.Handler
}

func NewServer() *http.Server {
	db := database.New()

	handler := handler.NewHandler(db)

	NewServer := &Server{
		port: 80,

		db:      db,
		handler: handler,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
