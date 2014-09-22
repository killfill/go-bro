package service

import (
	"go-bro/broker"

	"go-bro/service/postgresql"
)

func New(serviceType, connection string) broker.ServiceBroker {

	//TODO: reflection?
	switch serviceType {

	case "postgresql":
		return postgresql.New(connection)

	default:
		panic("Unknown Type: " + serviceType)
	}
}
