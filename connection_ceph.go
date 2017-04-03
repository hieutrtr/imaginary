package main

import (
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

func (c *CephConnection) writeToBlock(buf []byte) error {
	return ioutil.WriteFile(c.GetBlockPath(), buf, 0644)
}

// Execute purpose of openning connection
func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	var err error
	if !c.IsEnable() {
		return NewError("ceph: service is not supported", Unsupported)
	}
	c.BindRequest(r)

	if c.UseBlock {
		return c.writeToBlock(buf)
	}

	if !c.OnContext() {
		err = c.OpenContext()
		if err != nil {
			return err
		}
	}

	// Clear object before update original data
	if c.Attr == DATA {
		if err = c.DelObj(); err != nil {
			LoggerInfo.Println("WARNING: No object match with", c.CephObject, "to delete with reason", err)
		}
	}
	return c.SetAttr(buf)
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
