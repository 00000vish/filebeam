package main

import (
	"fmt"
	"os"

	"github.com/00000vish/filebeam/internals/config"
	"github.com/00000vish/filebeam/internals/connection"
	"github.com/00000vish/filebeam/internals/dirserver"
	"github.com/00000vish/filebeam/internals/endpoint"
	"github.com/00000vish/filebeam/internals/logger"
)

func main() {
	config := config.Parse()
	if !config.IsValid() {
		fmt.Println("todo")
	}

	logging := logger.Create(config)
	fmt.Println(logging)

	connection, error := connection.Create(&config)
	if error != nil {
		logger.Log(error.Error())
		os.Exit(1)
	}

	endpoint := endpoint.Endpoint{}
	error = endpoint.Create(connection)
	if error != nil {
		logger.Log(error.Error())
		os.Exit(1)
	}

	server := dirserver.DirServer{}
	server.Run(&config, &endpoint)

	select {}
}
