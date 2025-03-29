package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.RegisterRoutes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))

	}
}

func connect() (*amqp.Connection, error) {
	//TODO: connect to rabbitmq
	return nil, nil
}
