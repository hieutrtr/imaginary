package main

import (
	"fmt"
	"net/http"

	"github.com/noahdesu/go-ceph/rados"
)

const CephConnectionType ConnectionType = "ceph"

type CephConnection struct {
	Config *ConnectionConfig
	Conn   *rados.Conn
}

func NewCephConnection(config *ConnectionConfig) Connection {
	cephConn := &CephConnection{}
	cephConn.Config = config
	if config.EnableCeph {
		cephConn.Conn = MakeCephConnection(config)
	}
	return cephConn
}

func MakeCephConnection(config *ConnectionConfig) *rados.Conn {
	conn, err := rados.NewConn()
	if err != nil {
		exitWithError("rados connection fail: %s", err)
	}
	conn.ReadConfigFile(config.CephConfig)
	err = conn.Connect()
	if err != nil {
		exitWithError("rados connection fail: %s", err)
	}
	return conn
}

func (c *CephConnection) Matches(r *http.Request) bool {
	return r.Method == "POST" && r.URL.Query().Get("cpool") != "" && r.URL.Query().Get("coid") != "" && r.URL.Query().Get("cattr") != ""
}

func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	if c.Config.EnableCeph == false {
		return fmt.Errorf("Ceph is not enable")
	}
	ioctx, err := c.Conn.OpenIOContext(r.URL.Query().Get("cpool"))
	defer ioctx.Destroy()
	if err != nil {
		return err
	}
	err = ioctx.SetXattr(r.URL.Query().Get("coid"), r.URL.Query().Get("cattr"), buf)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
