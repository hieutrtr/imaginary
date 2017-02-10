package main

import "net/http"

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
	return r.Method == "POST" && r.URL.Query().Get("cpool") != "" && r.URL.Query().Get("coid") != ""
}

// Execute purpose of openning connection
func (c *CephConnection) Execute(r *http.Request, buf []byte) error {
	if !c.IsEnable() {
		return NewError("ceph: service is not supported", Unsupported)
	}
	c.BindRequest(r)
	err := c.OpenContext()
	if err != nil {
		return err
	}
	defer c.DestroyContext()
	err = c.SetData(buf)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterConnection(CephConnectionType, NewCephConnection)
}
