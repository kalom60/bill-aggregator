package main

import (
	"fmt"

	"github.com/kalom60/bill-aggregator/user-service/internal/server"
)

func main() {

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
