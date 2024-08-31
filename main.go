package main

import (
	"fmt"

	"github.com/dancankarani/palace/endpoints"
	"github.com/dancankarani/palace/model"
)

func main() {
	fmt.Println("hello world")
	model.MigrateDB()
	endpoints.CreateEndpoint()
}