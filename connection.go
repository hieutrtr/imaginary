package main

import "net/http"

type ConnectionType string
type ConnectionFactoryFunction func(*ConnectionConfig) Connection

type ConnectionConfig struct {
	PoolName   string
	EnableCeph bool
	CephConfig string
}

var connectionMap = make(map[ConnectionType]Connection)
var connectionFactoryMap = make(map[ConnectionType]ConnectionFactoryFunction)

type Connection interface {
	Matches(*http.Request) bool
	Execute(*http.Request) error
}

func RegisterConnection(connType ConnectionType, factory ConnectionFactoryFunction) {
	connectionFactoryMap[connType] = factory
}

func LoadConnections(o ServerOptions) {
	for name, factory := range connectionFactoryMap {
		connectionMap[name] = factory(&ConnectionConfig{
			PoolName:   PoolName,
			EnableCeph: o.EnableCeph,
			CephConfig: o.CephConfig,
		})
	}
}

func MatchConnection(req *http.Request) Connection {
	for _, conn := range connectionMap {
		if conn.Matches(req) {
			return conn
		}
	}
	return nil
}
