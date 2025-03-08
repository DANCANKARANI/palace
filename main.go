package main

import (
	"fmt"

	"github.com/dancankarani/palace/database"
	"github.com/dancankarani/palace/endpoints"
	"github.com/dancankarani/palace/model"
)

func main() {
	fmt.Println(".....")
	model.MigrateDB()
    endpoints.CreateEndpoint()
    database.ConnectDB()
}
