package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"bill-aggregation-service/internal/database"
	"bill-aggregation-service/internal/grpc/clients"
	"bill-aggregation-service/internal/handler"
)

type Server struct {
	port    int
	db      database.Service
	handler handler.Handler
}

func NewServer() *http.Server {
	db := database.New()

	gRPCAccountClient, err := clients.NewAccountClient("account-linking-service:50001")
	if err != nil {
		log.Fatalf("Failed to connect to provider service: %v", err)
	}
	defer gRPCAccountClient.Close()

	gRPCProviderClient, err := clients.NewProviderClient("utility-provider-service:50001")
	if err != nil {
		log.Fatalf("Failed to connect to provider service: %v", err)
	}
	defer gRPCProviderClient.Close()

	handler := handler.NewHandler(db, gRPCAccountClient, gRPCProviderClient)

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
