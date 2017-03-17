package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/noahdesu/go-ceph/rados"

	gorilla "github.com/gorilla/mux"
)

const (
	DATA               = "data"
	WIDTH              = "width"
	HEIGH              = "heigh"
	SIZE               = "size"
	IMAGE_MAX_BYTE     = 5242880
	CONNECTION_TIMEOUT = 10
	CTX_TIMEOUT        = 5
)

// Ceph main struct of ceph
type Ceph struct {
	Connection *rados.Conn
	Context    map[string]*rados.IOContext
	CephObject
	CephConfig
}

// CephConfig general config of ceph
type CephConfig struct {
	ConfigPath string
	Enable     bool
	UseBlock   bool
	BlockURL   string
}

// CephObject need to get ceph object
type CephObject struct {
	Pool string
	OID  string
}

// GetBlockPath build block storage path
func (c *Ceph) GetBlockPath() string {
	return fmt.Sprintf("/%s/%s/%s", c.BlockURL, c.Pool, c.OID)
}

// OnContext check if context in request is registered
func (c *Ceph) OnContext() bool {
	if c.Context[c.Pool] == nil {
		return false
	}
	return true
}

// IsEnable check ceph is served
func (c *Ceph) IsEnable() bool {
	return c.Enable
}

// SetData push data to ceph object
func (c *Ceph) SetData(buf []byte) error {
	errSignal := make(chan error, 1)

	go func() {
		errSignal <- c.Context[c.Pool].SetXattr(c.OID, DATA, buf)
	}()

	select {
	case <-time.After(time.Second * CTX_TIMEOUT):
		return NewError("Ceph Connection Timeout", 1)
	case err := <-errSignal:
		if err != nil {
			return err
		}
		return nil
	}
}

// GetData fetch object attribute DATA from ceph
func (c *Ceph) GetData() ([]byte, error) {
	errSignal := make(chan error, 1)
	lengSignal := make(chan int, 1)
	data := make([]byte, IMAGE_MAX_BYTE)
	go func() {
		leng, err := c.Context[c.Pool].GetXattr(c.OID, DATA, data)
		if err != nil {
			errSignal <- NewError("Data is not exists", NotFound)
		}
		// Remove any NULL characters from buffer
		if data == nil {
			errSignal <- NewError("Data is not exists", NotFound)
		}
		lengSignal <- leng
	}()

	select {
	case <-time.After(time.Second * CTX_TIMEOUT):
		return nil, NewError("Ceph Connection Timeout", 1)
	case err := <-errSignal:
		return nil, err
	case leng := <-lengSignal:
		buf := bytes.NewBuffer(make([]byte, 0, leng+1))
		io.Copy(buf, bytes.NewReader(data[:leng]))
		return buf.Bytes(), nil
	}
}

// DestroyContext when finish ceph jobs
func (c *Ceph) DestroyContext() {
	c.Context[c.Pool].Destroy()
}

// OpenContext provide context from existed pool on ceph cluster
// Should ask sysad to create pool on ceph first
func (c *Ceph) OpenContext() error {
	ioctx, err := c.Connection.OpenIOContext(c.Pool)
	if err != nil {
		return NewError("ceph: cannot open context of pool "+c.Pool, BadRequest)
	}
	if c.Context == nil {
		c.Context = make(map[string]*rados.IOContext)
	}
	c.Context[c.Pool] = ioctx
	return nil
}

// BindRequest Initialize CephObject need to get ceph object
func (c *Ceph) BindRequest(req *http.Request) {
	vars := gorilla.Vars(req)
	c.CephObject = CephObject{
		Pool: vars["service"],
		OID:  vars["oid"],
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
	return c.Connection.Connect()
}
