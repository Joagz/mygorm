package orm

import "fmt"

type Connector interface {
	init() error
	ModelOf(any, string) (Model, error)
}

type connector struct {
	datasourceConfig
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

	return &conn
}

func (c connector) init() error {
	switch c.Engine {
	case "mysql":
		initMySqlConnector(c)
	default:
		return fmt.Errorf("no connector provided for engine '%s'", c.Engine)
	}

	return nil
}

func initMySqlConnector(c connector) {

}
