package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Username string
	Password string
	Services []ServiceConfig
	Plans    map[string]PlanConfig
}

type ServiceConfig struct {
	Id, Type, Connection string
}

type PlanConfig struct {
	Concurrency int
}

func FromJson(path string) Config {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err.Error())
	}

	return config
}
