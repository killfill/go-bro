package main

import (
	"os"
	"go-bro/broker"
	"go-bro/config"

	"go-bro/service"

	"fmt"
)

func main() {

	addr := getAddr()
	config := config.FromJson("config.json")

	b := broker.New(config.Username, config.Password, config.Plans)

	//Setup known services
	for _, serviceConfig := range config.Services {

		fmt.Printf("Registering Service Broker %+v \n", serviceConfig)

		serviceBroker := service.New(serviceConfig.Type, serviceConfig.Connection)

		b.RegisterService(serviceConfig.Id, serviceBroker)
	}

	b.Listen(addr)

}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}

	return ":" + port
}
