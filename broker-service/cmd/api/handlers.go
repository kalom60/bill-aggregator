package main

import "github.com/gin-gonic/gin"

type RequestPayload struct {
	Action string `json:"action"`
}

func (app *Config) Broker(c *gin.Context) {
	c.JSON(200, gin.H{
		"Error":   false,
		"Message": "Hit the broker",
	})
}
