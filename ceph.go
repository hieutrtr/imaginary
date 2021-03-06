package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ceph/go-ceph/rados"

	gorilla "github.com/gorilla/mux"
)

const (
	DATA               = "data"
	CONNECTION_TIMEOUT = 10
	CTX_TIMEOUT        = 5
)

// Ceph main struct of ceph
type Ceph struct {
	Connection *rados.Conn
	Context    map[string]*rados.IOContext
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
func (c *Ceph) GetStat(obj *CephObject) (rados.ObjectStat, error) {
	if !c.OnContext(obj.Pool) {
		err := c.OpenContext(obj.Pool)
		if err != nil {
			return rados.ObjectStat{}, err
		}
	}
	return c.Context[obj.Pool].Stat(obj.OID)
}

// GetBlockPath build block storage path
func (c *Ceph) GetBlockPath(obj *CephObject) string {
	return fmt.Sprintf("/%s/%s/%s", c.BlockURL, obj.Pool, obj.OID)
}

// OnContext check if context in request is registered
func (c *Ceph) OnContext(Pool string) bool {
	if c.Context[Pool] == nil {
		return false
	}
	return true
}

// IsEnable check ceph is served
func (c *Ceph) IsEnable() bool {
	return c.Enable
}

// DelObj delete an object from ceph
func (c *Ceph) DelObj(obj *CephObject) error {
	if !c.OnContext(obj.Pool) {
		err := c.OpenContext(obj.Pool)
		if err != nil {
			return err
		}
	}
	errSignal := make(chan error, 1)
	go func() {
		errSignal <- c.Context[obj.Pool].Delete(obj.OID)
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
func (c *Ceph) SetAttr(obj *CephObject, buf []byte) error {
	if !c.OnContext(obj.Pool) {
		err := c.OpenContext(obj.Pool)
		if err != nil {
			return err
		}
	}

	errSignal := make(chan error, 1)
	go func() {
		if obj.Attr == "" {
			obj.Attr = DATA
		} else if obj.Attr != DATA {
			LoggerInfo.Println("cache Object's attribute", obj)
		}

		errSignal <- c.Context[obj.Pool].SetXattr(obj.OID, obj.Attr, buf)
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

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, *aMaxAllowedSize)
	},
}


// GetAttr fetch object attribute DATA from ceph
func (c *Ceph) GetAttr(obj *CephObject) ([]byte, error) {
	if !c.OnContext(obj.Pool) {
		err := c.OpenContext(obj.Pool)
		if err != nil {
			return nil, err
		}
	}
	
	data := bufPool.Get().([]byte)
	defer bufPool.Put(data)
	
	if obj.Attr == "" {
		obj.Attr = DATA
	}
	length, err := c.Context[obj.Pool].GetXattr(obj.OID, obj.Attr, data)
	if err != nil {
	    return nil, err
	}
		
	buf := make([]byte, length)
	copy(buf, data)
	return buf, nil
}

// DestroyContext when finish ceph jobs
// func (c *Ceph) DestroyContext() {
// 	c.Context[c.Pool].Destroy()
// }

// OpenContext provide context from existed pool on ceph cluster
// Should ask sysad to create pool on ceph first
func (c *Ceph) OpenContext(Pool string) error {
	ioctx, err := c.Connection.OpenIOContext(Pool)
	if err != nil {
		return NewError("ceph: cannot open context of pool "+Pool, BadRequest)
	}
	if c.Context == nil {
		c.Context = make(map[string]*rados.IOContext)
	}
	c.Context[Pool] = ioctx
	return nil
}

// BindRequest Initialize CephObject need to get ceph object
func BindRequest(req *http.Request) *CephObject {
	vars := gorilla.Vars(req)
	attr := getCacheAttr(req.URL.Path, req.URL.RawQuery)
	return &CephObject{
		Pool: vars["service"],
		OID:  vars["oid"],
		Attr: attr,
	}
}

// BindObject Initialize CephObject to get ceph object by attribute
func BindObject(vars map[string]string) *CephObject {
	return &CephObject{
		Pool: vars["service"],
		OID:  vars["oid"],
		Attr: DATA,
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
