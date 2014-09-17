package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"sql-broker/broker"
	"sql-broker/service/postgres"
)

func getDSN() string {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://postgres@db-server:5432/?sslmode=disable&connect_timeout=10"
	}
	return dsn
}

func LoadConfig(file string) (config broker.Config, err error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return
	}
	return
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return port
}

func main() {

	service, err := postgres.New(getDSN())
	if err != nil {
		fmt.Println("Could not connect", err.Error())
		os.Exit(1)
	}

	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Println("Could not load config:", config)
		os.Exit(2)
	}
	// fmt.Println("CACHATE", config)

	port := getPort()
	fmt.Println("Starting on port", port)
	broker.Start(config, service, ":"+port)

}
