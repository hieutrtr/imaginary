package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/noahdesu/go-ceph/rados"

	gorilla "github.com/gorilla/mux"
)

const (
	DATA               = "data"
	WIDTH              = "width"
	HEIGH              = "heigh"
	SIZE               = "size"
	IMAGE_MAX_BYTE     = 20971520
	CONNECTION_TIMEOUT = 10
	CTX_TIMEOUT        = 5
)

var cephAttributes = []string{
	"thumbnail",
	"watermark",
}

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
	Attr string
}

// GetStat of object from ceph
// TODO: Need improve with handling connection timeout
func (c *Ceph) GetStat() (rados.ObjectStat, error) {
	return c.Context[c.Pool].Stat(c.OID)
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

// DelObj delete an object from ceph
func (c *Ceph) DelObj() error {
	errSignal := make(chan error, 1)
	go func() {
		errSignal <- c.Context[c.Pool].Delete(c.OID)
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

// SetAttr push attribute to ceph object
func (c *Ceph) SetAttr(buf []byte) error {
	errSignal := make(chan error, 1)
	go func() {
		if c.Attr == "" {
			c.Attr = DATA
		} else if c.Attr != DATA {
			LoggerInfo.Println("cache Object's attribute", c.CephObject)
		}
		errSignal <- c.Context[c.Pool].SetXattr(c.OID, c.Attr, buf)
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

// GetAttr fetch object attribute DATA from ceph
func (c *Ceph) GetAttr() ([]byte, error) {
	errSignal := make(chan error, 1)
	lengSignal := make(chan int, 1)
	data := make([]byte, IMAGE_MAX_BYTE)
	go func() {
		if c.Attr == "" {
			c.Attr = DATA
		}
		leng, err := c.Context[c.Pool].GetXattr(c.OID, c.Attr, data)
		if err != nil {
			errSignal <- NewError(err.Error(), NotFound)
		}
		// Remove any NULL characters from buffer
		if data == nil {
			errSignal <- NewError("No data", NotFound)
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
	attr := getCacheAttr(req)
	c.CephObject = CephObject{
		Pool: vars["service"],
		OID:  vars["oid"],
		Attr: attr,
	}
}

func getCacheAttr(req *http.Request) string {
	parts := strings.Split(req.URL.Path, "/")
	if parts[1] != "upload" {
		for _, a := range cephAttributes {
			if a == parts[len(parts)-1] {
				attr := fmt.Sprintf("%s_%s", a, req.URL.RawQuery)
				return attr
			}
		}
	}
	return DATA
}

// BindObject Initialize CephObject to get ceph object by attribute
func (c *Ceph) BindObject(vars map[string]string) {
	c.CephObject = CephObject{
		Pool: vars["service"],
		OID:  vars["oid"],
		Attr: vars["attr"],
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
