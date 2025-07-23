package orm

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Connector interface {
	Open() error
	Close()
	ModelOf(any, string) (Model, error)
}

type connector struct {
	datasourceConfig
}

type mySqlConnector struct {
	connector
	*sql.DB
}

func NewConnector(engine, username, password, dbname, hostname string) Connector {
	cfg := datasourceConfig{
		Engine:   engine,
		username: username,
		password: password,
		Schema:   dbname,
		Hostname: hostname,
	}

	conn := connector{
		datasourceConfig: cfg,
	}

	switch engine {
	case "mysql":
		fmt.Printf("returning mySqlConnector\n")
		return &mySqlConnector{
			connector: conn,
			DB:        nil,
		}
	}

	return &conn
}

func (c connector) Close() {
	panic("please select a sql engine")
}

func (c connector) Open() error {
	panic("please select a sql engine")
}

func (c mySqlConnector) Close() {
	if c.DB == nil {
		panic("please use Open() before Close()")
	}
	err := c.DB.Close()
	if err != nil {
		fmt.Printf("Failed to close mySqlConnector: '%s'\n", err.Error())
	}
}

func makeMySqlConnectionString(c mySqlConnector) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?collation=utf8mb4_unicode_ci&parseTime=true", c.username, c.password, c.Hostname, c.Schema)
}

func (c *mySqlConnector) Open() error {
	fmt.Println("mysql.Open() called") // Debug
	db, err := sql.Open("mysql", makeMySqlConnectionString(*c))
	if err != nil {
		return err	
	}

	c.DB = db
	return nil
}
