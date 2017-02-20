package main

import (
	"bytes"
	"net/http"

	"github.com/noahdesu/go-ceph/rados"
)

const (
	DATA           = "data"
	WIDTH          = "width"
	HEIGH          = "heigh"
	SIZE           = "size"
	IMAGE_MAX_BYTE = 5242880
)

// Ceph main struct of ceph
type Ceph struct {
	Connection *rados.Conn
	Context    *rados.IOContext
	CephObject
	CephConfig
}

// CephConfig general config of ceph
type CephConfig struct {
	ConfigPath string
	Enable     bool
}

// CephObject need to get ceph object
type CephObject struct {
	Pool string
	OID  string
}

// IsEnable check ceph is served
func (c *Ceph) IsEnable() bool {
	return c.Enable
}

// SetData push data to ceph object
func (c *Ceph) SetData(buf []byte) error {
	if err := c.Context.SetXattr(c.OID, DATA, buf); err != nil {
		return err
	}
	return nil
}

// GetData fetch object attribute DATA from ceph
func (c *Ceph) GetData() ([]byte, error) {
	buf := make([]byte, IMAGE_MAX_BYTE)
	if _, err := c.Context.GetXattr(c.OID, DATA, buf); err != nil {
		return nil, NewError("Data is not exists", NotFound)
	}
	// Remove any NULL characters from buffer
	buf = bytes.Trim(buf, "\x00")
	return buf, nil
}

// DestroyContext when finish ceph jobs
func (c *Ceph) DestroyContext() {
	c.Context.Destroy()
}

// OpenContext provide context from existed pool on ceph cluster
// Should ask sysad to create pool on ceph first
func (c *Ceph) OpenContext() error {
	ioctx, err := c.Connection.OpenIOContext(c.Pool)
	c.Context = ioctx
	if err != nil {
		return NewError("ceph: cannot open context of pool "+c.Pool, BadRequest)
	}
	return nil
}

// BindRequest Initialize CephObject need to get ceph object
func (c *Ceph) BindRequest(req *http.Request) {
	c.CephObject = CephObject{
		Pool: req.URL.Query().Get("cpool"),
		OID:  req.URL.Query().Get("coid"),
	}
}

// Connect do connect steps with ceph server by config path
// ConfigPath : /etc/ceph/ceph.conf by default
func (c *Ceph) Connect() error {
	c.Connection, _ = rados.NewConn()
	err := c.Connection.ReadConfigFile(c.ConfigPath)
	if err != nil {
		return NewError("ceph: fail to read config "+c.ConfigPath, NotFound)
	}
	err = c.Connection.Connect()
	if err != nil {
		return NewError("ceph: fail to connect with config "+c.ConfigPath, NotFound)
	}
	return nil

}
