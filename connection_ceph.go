package main

import (
	"net/http"

	gorilla "github.com/gorilla/mux"
)

// CephConnectionType name of regiter connection
const CephConnectionType ConnectionType = "ceph"

// CephConnection to register for Ceph
type CephConnection struct {
	Ceph
}

// NewCephConnection create new ceph connection
func NewCephConnection(config *ConnectionConfig) Connection {
	cc := &CephConnection{
		Ceph: Ceph{
			CephConfig: CephConfig{
				ConfigPath: config.CephConfig,
				Enable:     config.EnableCeph,
			},
		},
	}
	if config.EnableCeph {
		MakeConnection(cc)
	}
	return cc
}

func (c *CephConnection) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "POST" && vars["cpool"] != "" && vars["coid"] != ""
}

// Execute purpose of openning connection
func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	if !c.IsEnable() {
		return NewError("ceph: service is not supported", Unsupported)
	}
	c.BindRequest(r)
	if c.Context == nil {
		err := c.OpenContext()
		if err != nil {
			return err
		}
	}
	return c.SetData(buf)
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
