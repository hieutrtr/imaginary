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

func NewCephConnection(config *ConnectionConfig) *CephConnection {
	cephConn := &CephConnection{}
	ceph.Config = config
	ceph.Conn = MakeConnection(config)
	return cephConn
}

func MakeConnection(config *ConnectionConfig) *rados.Conn {
	return nil
}

func (c *CephConnection) Matches(r *http.Request) bool {
	return r.Method == "POST" && r.URL.Query().Get("cns") != "" && r.URL.Query().Get("cid") != ""
}

func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	fmt.Println(buf)
	return nil
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
