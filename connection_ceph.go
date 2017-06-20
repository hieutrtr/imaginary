package main

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	if config.EnableCeph && !config.UseCephBlock {
		cc := &CephConnection{
			Ceph: Ceph{
				CephConfig: CephConfig{
					ConfigPath: config.CephConfig,
					Enable:     config.EnableCeph,
					UseBlock:   config.UseCephBlock,
					BlockURL:   config.CephBlockURL,
				},
			},
		}
		err := MakeConnection(cc)
		if err != nil {
			exitWithError("Ceph connection was fail %s with config: %s", fmt.Sprint(err), config.CephConfig)
		}
		return cc
	}
	return nil
}

func (c *CephConnection) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "POST" && vars["service"] != "" && vars["oid"] != ""
}

// Execute purpose of openning connection
func (c *CephConnection) Execute(req *http.Request, buf []byte) error {
	// var err error
	if !c.IsEnable() {
		return errors.New("ceph: service is not supported")
	}

	vars := gorilla.Vars(req)

	if c.UseBlock {
		return ioutil.WriteFile(c.GetBlockPath(BindObject(vars)), buf, 0644)
	}

	return c.SetAttr(BindObject(vars), buf)
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
