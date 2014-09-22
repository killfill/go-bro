package service

import (
	"go-bro/broker"
	"go-bro/config"
	"go-bro/service/mysql"
	"go-bro/service/postgresql"
)

func New(conf config.ServiceConfig) broker.ServiceBroker {

	//TODO: reflection?
	switch conf.Type {

	case "postgresql":
		return postgresql.New(conf)

	case "mysql":
		return mysql.New(conf)

	default:
		panic("Unknown Type: " + conf.Type)
	}
}
