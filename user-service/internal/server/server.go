package server

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"user-service/internal/database"
	"user-service/internal/handler"
)

type Server struct {
	port int

	db      database.Service
	handler handler.Handler
}

func NewServer() *http.Server {
	db := database.New()

	userHandler := handler.NewHandler(db)

	NewServer := &Server{
		port:    80,
		db:      db,
		handler: userHandler,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
