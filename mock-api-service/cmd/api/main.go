package main

import (
	"fmt"
	"net/http"
)

const port = "80"

type Config struct{}

func main() {
	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.RegisterRoutes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))

	}
}
