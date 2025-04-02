package server

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/database"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/grpc"
	"github.com/kalom60/bill-aggregator/utility-provider-service/internal/handler"
)

type Server struct {
	port    int
	db      database.Service
	handler handler.Handler
}

func NewServer() *Server {
	db := database.New()
	handler := handler.NewHandler(db)

	return &Server{
		port:    80,
		db:      db,
		handler: handler,
	}

}

func (s *Server) ListenAndServe() error {
	go grpc.StartGRPCListen(s.db)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.RegisterRoutes(),
	}

	fmt.Printf("HTTP server running on port %d\n", s.port)
	return server.ListenAndServe()
}
