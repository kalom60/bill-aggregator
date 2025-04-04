package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kalom60/bill-aggregator/broker/pkg/middlewares"
)

func (app *Config) RegisterRoutes() http.Handler {
	r := gin.Default()

	middlewares.InitRedis()

	r.Use(middlewares.RateLimiterMiddleware(10, time.Minute))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.POST("/", app.Broker)
	r.POST("/handle", app.HandleSubmission)

	return r
}
