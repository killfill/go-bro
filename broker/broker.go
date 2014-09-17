package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

var (
	service ServiceBroker
	config  Config
)

func Start(c Config, s ServiceBroker, port string) {

	service = s
	config = c

	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(port, nil)
}

// func usersHandler(w http.ResponseWriter, r *http.Request) {

// 	users, _ := service.GetUsers()
// 	js, err := json.Marshal(users)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)

// }

func idsFromPath(path string) (serviceID string, bindID string) {

	x := strings.Split(path, "/")
	serviceID = x[1]
	if len(x) == 4 {
		bindID = x[3]
	}

	return
}

var empty struct{} = struct{}{}

func serviceHandler(r *http.Request) (int, interface{}) {

	serviceID, bindingID := idsFromPath(r.URL.Path[len("/v2/"):])

	decoder := json.NewDecoder(r.Body)

	switch {

	//CREATE BINDING
	case bindingID != "" && r.Method == "PUT":

		var req BindRequest

		if err := decoder.Decode(&req); err != nil {
			return http.StatusBadRequest, err
		}

		resp, err := service.Bind(serviceID, bindingID, req)
		if err != nil {
			return http.StatusConflict, err
		}

		return http.StatusCreated, resp

	case bindingID != "" && r.Method == "DELETE":
		if err := service.Unbind(serviceID, bindingID); err != nil {
			return http.StatusGone, err
		}

		return http.StatusOK, empty

	//CREATE SERVICE
	case r.Method == "PUT":

		var req ServiceRequest

		if err := decoder.Decode(&req); err != nil {
			return http.StatusBadRequest, err
		}

		if err := service.Create(serviceID, req); err != nil {
			return http.StatusConflict, err
		}

		return http.StatusCreated, empty

	//DELETE SERVICE
	case r.Method == "DELETE":

		if err := service.Destroy(serviceID); err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, empty

	default:
		return 404, empty
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {

	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println("\n---------------------------")
	fmt.Println(string(dump))

	switch {

	case strings.HasPrefix(r.URL.Path, "/v2/catalog"):

		http.ServeFile(w, r, "catalog.json")

	case strings.HasPrefix(r.URL.Path, "/v2/service_instances/"):

		status, body := serviceHandler(r)

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
