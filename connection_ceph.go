package main

import (
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
	if config.EnableCeph && !config.UseCephBlock {
		err := MakeConnection(cc)
		if err != nil {
			exitWithError("Ceph connection was fail with config: %s", config.CephConfig)
		}
	}
	return cc
}

func (c *CephConnection) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "POST" && vars["cpool"] != "" && vars["coid"] != ""
}

func (c *CephConnection) writeToBlock(buf []byte) error {
	return ioutil.WriteFile(c.GetBlockPath(), buf, 0644)
}

// Execute purpose of openning connection
func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	if !c.IsEnable() {
		return NewError("ceph: service is not supported", Unsupported)
	}
	c.BindRequest(r)

	if c.UseBlock {
		return c.writeToBlock(buf)
	}

	if !c.OnContext() {
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
