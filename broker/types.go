package broker

import "go-bro/config"

type ServiceRequest struct {
	ServiceID      string `json:"service_id"`
	PlanID         string `json:"plan_id"`
	OrganizationID string `json:"organization_guid"`
	SpaceID        string `json:"space_guid"`
}

type BindRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
	AppID     string `json:"app_guid"`
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

type ServiceBroker interface {
	Create(string, ServiceRequest, config.PlanConfig) error
	Destroy(string) error

	Bind(string, string, BindRequest) (BindResponse, error)
	Unbind(string, string) error
}
