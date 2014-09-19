package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"

	"sql-broker/broker"

	"fmt"
	"net/url"
	"strconv"
	"strings"

	"crypto/rand"
)

type PostgresService struct {
	db     sql.DB
	dsnTpl broker.BindResponse
}

type User struct {
	Usename   string
	Id        int
	CanCreate bool
}

func UserCredentialTpl(dsn string) broker.BindResponse {
	u, _ := url.Parse(dsn)

	s := strings.Split(u.Host, ":")

	//What a mess
	var host string
	var port int
	if len(s) > 1 {
		host = s[0]
		port, _ = strconv.Atoi(s[1])
	} else {
		host, port = s[0], 5432
	}

	return broker.BindResponse{Credentials: broker.BindCredentials{Host: host, Port: port}}
}

//"Static method"
func New(dsn string) (PostgresService, error) {

	conn, err := sql.Open("postgres", dsn)
	defer conn.Close()

	s := PostgresService{db: *conn}

	if err != nil {
		return s, err
	}

	//Lets reuse the dsn as a template to build the user credentials object when binding.
	s.dsnTpl = UserCredentialTpl(dsn)

	return s, conn.Ping()
}

func (s PostgresService) Create(serviceInstance string, req broker.ServiceRequest, limit broker.Limit) error {

	sql := fmt.Sprintf("CREATE DATABASE \"%s\" CONNECTION LIMIT %d", serviceInstance, limit.Concurrency)

	_, err := s.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresService) Destroy(serviceInstance string) error {

	_, err := s.db.Exec(fmt.Sprintf("DROP DATABASE \"%s\"", serviceInstance))
	if err != nil {
		return err
	}

	return nil
}

func (s PostgresService) Bind(serviceInstance string, bindID string, req broker.BindRequest) (broker.BindResponse, error) {

	//TODO: Make a copy of this thing
	s.dsnTpl.Credentials.Username = bindID
	s.dsnTpl.Credentials.Password = rand_str(15)
	s.dsnTpl.Credentials.Database = serviceInstance
	s.dsnTpl.Credentials.Uri = s.dsnTpl.Credentials.String()

	sql := fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", s.dsnTpl.Credentials.Username, s.dsnTpl.Credentials.Password)
	_, err := s.db.Exec(sql)
	if err != nil {
		return s.dsnTpl, err
	}

	sql = fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE \"%s\" to \"%s\"", s.dsnTpl.Credentials.Database, s.dsnTpl.Credentials.Username)
	_, err = s.db.Exec(sql)
	if err != nil {
		return s.dsnTpl, err
	}

	return s.dsnTpl, nil
}

func (s PostgresService) Unbind(serviceInstance string, bindID string) error {

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

//Taken from http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go
func rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
