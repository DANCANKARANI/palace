package main

import (
	"fmt"

	"github.com/dancankarani/palace/database"
	"github.com/dancankarani/palace/endpoints"
)

func main() {
	fmt.Println(".....")
    endpoints.CreateEndpoint()
    database.ConnectDB()
}
