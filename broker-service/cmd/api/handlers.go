package main

import "github.com/gin-gonic/gin"

type RequestPayload struct {
	Action string `json:"action"`
}

func (app *Config) Broker(c *gin.Context) {
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Hit the broker",
	})
}
