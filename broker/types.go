package broker

import (
	"fmt"
)

type Config struct {
	Username, Password string
	// Limits             []Limit
}

type Limit struct {
	Size, Concurrency int
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
	Create(string, ServiceRequest) error
	Destroy(string) error

	Bind(string, string, BindRequest) (BindResponse, error)
	Unbind(string, string) error
}

// // ESTOS PARECE QUE YA NO LOS OCUPO!!
// type Catalog struct {
// 	Services []Service
// }

// type Service struct {
// 	Id, Name, Description string
// 	Bindable              bool
// 	Tags, Requires        []string
// 	Metadata              ServiceMetadata
// 	Plans                 []Plan
// 	DashboardClient       ServiceDashboard
// }

// type ServiceMetadata struct {
// 	DisplayName, ImageURL, LongDescription, ProviderDisplayName, DocumentationUrl, SupportURL string
// }

// type Plan struct {
// 	Id, Name, Description string
// 	Free                  bool
// 	Metadata              PlanMetadata
// }

// type PlanMetadata struct {
// 	Bullets     []string
// 	Costs       []map[string]interface{}
// 	DisplayName string
// }

// type ServiceDashboard struct {
// 	Id, Secret, RedirectUri string
// }
