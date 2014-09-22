package postgresql

import (
	"database/sql"
	_ "github.com/lib/pq"

	"sql-broker/broker"
	"sql-broker/config"
	"sql-broker/service/common"

	"fmt"
)

type PostgresService struct {
	db   sql.DB
	host string
	port int
}

func New(dsn string) *PostgresService {

	conn, err := sql.Open("postgres", dsn)
	defer conn.Close()

	if err != nil {
		panic(err.Error())
	}

	s := PostgresService{db: *conn}

	//Save the endpoint of the service, for later use in Bind()
	s.host, s.port = utils.GetAddressFromURL(dsn)

	return &s
}

func (s *PostgresService) Create(serviceInstance string, req broker.ServiceRequest, planConfig config.PlanConfig) (err error) {

	sql := fmt.Sprintf("CREATE DATABASE \"%s\" CONNECTION LIMIT %d", serviceInstance, planConfig.Concurrency)
	_, err = s.db.Exec(sql)
	return
}

func (s *PostgresService) Destroy(serviceInstance string) (err error) {

	_, err = s.db.Exec(fmt.Sprintf("DROP DATABASE \"%s\"", serviceInstance))
	return
}

func (s *PostgresService) Bind(serviceInstance string, bindID string, req broker.BindRequest) (resp broker.BindResponse, err error) {

	cred := broker.BindCredentials{Host: s.host, Port: s.port}
	cred.Username = bindID
	cred.Password = utils.Rand_str(15)
	cred.Database = serviceInstance
	cred.Uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cred.Username, cred.Password, cred.Host, cred.Port, cred.Database)

	resp = broker.BindResponse{Credentials: cred}

	sql := fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", cred.Username, cred.Password)
	_, err = s.db.Exec(sql)
	if err != nil {
		return
	}

	sql = fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE \"%s\" to \"%s\"", cred.Database, cred.Username)
	_, err = s.db.Exec(sql)
	return
}

func (s *PostgresService) Unbind(serviceInstance string, bindID string) error {

	sql := fmt.Sprintf("DROP OWNED BY \"%s\"", bindID)
	_, err := s.db.Exec(sql)
	if err != nil {
		return err
	}

	sql = fmt.Sprintf("DROP USER \"%s\"", bindID)
	_, err = s.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
