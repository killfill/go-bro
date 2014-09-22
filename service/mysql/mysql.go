package mysql

import (
	// "crypto/md5"
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"go-bro/broker"
	"go-bro/config"
	"go-bro/service/common"
)

type MysqlService struct {
	db   sql.DB
	conf config.ServiceConfig
}

func New(conf config.ServiceConfig) *MysqlService {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", conf.Username, conf.Password, conf.Host, conf.Port)
	conn, err := sql.Open("mysql", dsn)
	defer conn.Close()

	if err != nil {
		panic(err.Error())
	}

	// connectError := conn.Ping()
	// if connectError != nil {
	// 	panic("Could not connect: " + connectError.Error())
	// }

	return &MysqlService{db: *conn, conf: conf}
}

func (s MysqlService) Create(serviceInstance string, req broker.ServiceRequest, planConfig config.PlanConfig) error {

	// connection.execute("CREATE DATABASE IF NOT EXISTS #{connection.quote_table_name(database_name)}")
	sql := fmt.Sprintf("CREATE DATABASE `%s`", serviceInstance)
	_, err := s.db.Exec(sql)
	return err
}

func (s MysqlService) Destroy(serviceInstance string) error {

	_, err := s.db.Exec(fmt.Sprintf("DROP DATABASE `%s`", serviceInstance))
	return err
}

func (s MysqlService) Bind(serviceInstance string, bindID string, req broker.BindRequest, planConfig config.PlanConfig) (resp broker.BindResponse, err error) {

	cred := broker.BindCredentials{Host: s.conf.Host, Port: s.conf.Port}

	// md5Sum := md5.Sum([]byte(bindID))
	// cred.Username = fmt.Sprintf("%x", md5Sum[0:8]) //Mysql supports username <= 16 chars.

	cred.Username = bindID[0:16] //Mysql supports username <= 16 chars.
	cred.Password = utils.Rand_str(16)
	cred.Database = serviceInstance
	cred.Uri = fmt.Sprintf("mysql://%s:%s@%s:%d/%s", cred.Username, cred.Password, cred.Host, cred.Port, cred.Database)

	resp = broker.BindResponse{Credentials: cred}

	sql := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", cred.Username, cred.Password)
	_, err = s.db.Exec(sql)
	if err != nil {
		return
	}

	sql = fmt.Sprintf("GRANT ALL ON `%s`.* TO '%s'@'%%' WITH MAX_USER_CONNECTIONS %d", cred.Database, cred.Username, planConfig.Concurrency)
	_, err = s.db.Exec(sql)
	if err != nil {
		return
	}

	_, err = s.db.Exec("FLUSH PRIVILEGES")
	return

	return broker.BindResponse{}, nil
}

func (s MysqlService) Unbind(serviceInstance string, bindID string) error {

	sql := fmt.Sprintf("DROP USER '%s'", bindID[0:16])
	_, err := s.db.Exec(sql)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("FLUSH PRIVILEGES")
	return err
}
