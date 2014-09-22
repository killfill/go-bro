package service

import (
	"sql-broker/broker"

	"sql-broker/service/postgresql"
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
