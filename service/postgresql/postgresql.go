package postgresql

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"

	"go-bro/broker"
	"go-bro/config"
	"go-bro/service/common"
)

type PostgresService struct {
	db   sql.DB
	conf config.ServiceConfig
}

func New(conf config.ServiceConfig) *PostgresService {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=disable&connect_timeout=10", conf.Username, conf.Password, conf.Host, conf.Port)
	conn, err := sql.Open("postgres", dsn)
	defer conn.Close()

	if err != nil {
		panic(err.Error())
	}

	// connectError := conn.Ping()
	// if connectError != nil {
	// 	panic("Could not connect: " + connectError.Error())
	// }

	return &PostgresService{db: *conn, conf: conf}
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

	cred := broker.BindCredentials{Host: s.conf.Host, Port: s.conf.Port}
	cred.Username = bindID
	cred.Password = utils.Rand_str(16)
	cred.Database = serviceInstance
	cred.Uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cred.Username, cred.Password, cred.Host, cred.Port, cred.Database)

	resp = broker.BindResponse{Credentials: cred}

	sql := fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", cred.Username, cred.Password)
	_, err = s.db.Exec(sql)
	if err != nil {
		return
	}

	sql = fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\"", cred.Database, cred.Username)
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
