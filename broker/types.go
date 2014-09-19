package broker

import (
	"fmt"
)

type Config struct {
	Username, Password string
	Limits             map[string]Limit
}

type Limit struct {
	Concurrency int
}

type ServiceRequest struct {
	Service      string `json:"service_id"`
	Plan         string `json:"plan_id"`
	Organization string `json:"organization_guid"`
	Space        string `json:"space_guid"`
}

type BindRequest struct {
	Service string `json:"service_id"`
	Plan    string `json:"plan_id"`
	App     string `json:"app_guid"`
}

type BindResponse struct {
	Credentials BindCredentials `json:"credentials"`
}

type BindCredentials struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func (b *BindCredentials) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", b.Username, b.Password, b.Host, b.Port, b.Database)
}

type ServiceBroker interface {
	Create(string, ServiceRequest, Limit) error
	Destroy(string) error

	Bind(string, string, BindRequest) (BindResponse, error)
	Unbind(string, string) error
}
