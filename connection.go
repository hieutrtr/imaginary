package main

import (
	"net/http"
	"time"
)

type ConnectionType string
type ConnectionFactoryFunction func(*ConnectionConfig) Connection

type Connector interface {
	Connect() error
}

type ConnectionConfig struct {
	EnableS3     bool
	EnableCeph   bool
	CephConfig   string
	UseCephBlock bool
	CephBlockURL string
}

var connectionMap = make(map[ConnectionType]Connection)
var connectionFactoryMap = make(map[ConnectionType]ConnectionFactoryFunction)

type Connection interface {
	Matches(*http.Request) bool
	Execute(*http.Request, []byte) error
}

// MakeConnection connect by handler of connection
func MakeConnection(conn Connector) error {
	connSignal := make(chan error, 1)

	go func() {
		connSignal <- conn.Connect()
	}()

	select {
	case <-time.After(time.Second * CONNECTION_TIMEOUT):
		return NewError("Connection Timeout", 1)
	case err := <-connSignal:
		return err
	}
}

func RegisterConnection(connType ConnectionType, factory ConnectionFactoryFunction) {
	connectionFactoryMap[connType] = factory
}

func LoadConnections(o ServerOptions) {
	for name, factory := range connectionFactoryMap {
		connectionMap[name] = factory(&ConnectionConfig{
			EnableS3:     o.EnableS3,
			EnableCeph:   o.EnableCeph,
			CephConfig:   o.CephConfig,
			UseCephBlock: o.UseCephBlock,
			CephBlockURL: o.CephBlockURL,
		})
	}
}

func MatchConnection(req *http.Request) Connection {
	for _, conn := range connectionMap {
		if conn != nil && conn.Matches(req) {
			return conn
		}
	}
	return nil
}
