package broker

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"go-bro/config"
)

type Broker struct {
	username string
	password string
	plans    map[string]config.PlanConfig
	services map[string]ServiceBroker
}

func New(user, pass string, plans map[string]config.PlanConfig) Broker {
	return Broker{username: user, password: pass, plans: plans, services: map[string]ServiceBroker{}}
}

func (b *Broker) RegisterService(id string, service ServiceBroker) {
	b.services[id] = service
}

func (b *Broker) Listen(addr string) {
	fmt.Println("Starting on", addr)

	http.HandleFunc("/", b.httpHandler)
	http.ListenAndServe(addr, nil)
}

var empty struct{} = struct{}{}

func (b *Broker) serviceHandler(r *http.Request) (int, interface{}) {

	serviceInstance, bindingID := idsFromPath(r.URL.Path[len("/v2/"):])

	decoder := json.NewDecoder(r.Body)

	switch {

	//CREATE BINDING
	case bindingID != "" && r.Method == "PUT":

		var req BindRequest

		if err := decoder.Decode(&req); err != nil {
			return http.StatusBadRequest, err
		}

		service, gotService := b.services[req.ServiceID]
		if !gotService {
			return http.StatusBadRequest, errors.New("Unknown service: " + req.ServiceID)
		}

		resp, err := service.Bind(serviceInstance, bindingID, req)
		if err != nil {
			return http.StatusConflict, err
		}

		return http.StatusCreated, resp

	case bindingID != "" && r.Method == "DELETE":

		serviceID := r.URL.Query().Get("service_id")
		service, gotService := b.services[serviceID]
		if !gotService {
			return http.StatusBadRequest, errors.New("Unknown service: " + serviceID)
		}

		if err := service.Unbind(serviceInstance, bindingID); err != nil {
			return http.StatusGone, err
		}

		return http.StatusOK, empty

	//CREATE SERVICE
	case r.Method == "PUT":

		var req ServiceRequest

		if err := decoder.Decode(&req); err != nil {
			return http.StatusBadRequest, err
		}

		planConfig, gotPlan := b.plans[req.PlanID]
		if !gotPlan {
			return http.StatusBadRequest, errors.New("Unknown plan")
		}

		service, gotService := b.services[req.ServiceID]
		if !gotService {
			return http.StatusBadRequest, errors.New("Unknown service: " + req.ServiceID)
		}

		if err := service.Create(serviceInstance, req, planConfig); err != nil {
			return http.StatusConflict, err
		}

		return http.StatusCreated, empty

	//DELETE SERVICE
	case r.Method == "DELETE":

		serviceID := r.URL.Query().Get("service_id")
		service, gotService := b.services[serviceID]
		if !gotService {
			return http.StatusBadRequest, errors.New("Unknown service: " + serviceID)
		}

		if err := service.Destroy(serviceInstance); err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, empty

	default:
		return 404, empty
	}
}

func (b *Broker) httpHandler(w http.ResponseWriter, r *http.Request) {

	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println("\n---------------------------")
	fmt.Println(string(dump))

	auth := b.validCredentials(r.Header["Authorization"])
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch {

	case strings.HasPrefix(r.URL.Path, "/v2/catalog"):

		http.ServeFile(w, r, "catalog.json")

	case strings.HasPrefix(r.URL.Path, "/v2/service_instances/"):

		status, body := b.serviceHandler(r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		//If its an error, put the Message inside {'description': the_msg}
		err, found := body.(error)
		if found {
			body = map[string]string{"description": err.Error()}
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println("Could jsonify the body:", err)
		}

	default:
		http.NotFound(w, r)
	}

}

func (b *Broker) validCredentials(authHeader []string) bool {

	if len(authHeader) < 1 {
		return false
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 {
		return false
	}

	return parts[1] == base64.StdEncoding.EncodeToString([]byte(b.username+":"+b.password))
}

func idsFromPath(path string) (serviceID string, bindID string) {

	x := strings.Split(path, "/")
	serviceID = x[1]
	if len(x) == 4 {
		bindID = x[3]
	}

	return
}
