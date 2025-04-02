package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/database"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/grpc"
	"github.com/kalom60/bill-aggregator/account-linking-service/internal/handler"
)

type Server struct {
	port    int
	db      database.Service
	handler handler.Handler
}

func NewServer() *http.Server {
	db := database.New()

	gRPCClient, err := grpc.NewProviderClient("utility-provider-service:50001")
	if err != nil {
		log.Fatalf("Failed to connect to provider service: %v", err)
	}
	defer gRPCClient.Close()

	handler := handler.NewHandler(db, gRPCClient)

	NewServer := &Server{
		port:    80,
		db:      db,
		handler: handler,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
