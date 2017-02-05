package main

import (
	"net/http"

	"github.com/noahdesu/go-ceph/rados"
)

const ImageSourceTypeCeph ImageSourceType = "ceph"
const PoolName string = "media"

type CephObject struct {
	NameSpace string
	OID       string
}

type CephImageSource struct {
	Config     *SourceConfig
	Connection *rados.Conn
}

func NewCephImageSource(config *SourceConfig) ImageSource {
	CISource := &CephImageSource{}
	if config.EnableCeph {
		CISource.Config = config
		CISource.Connection = MakeConnection()
	}
	return CISource
}

func MakeConnection() *rados.Conn {
	conn, err := rados.NewConn()
	if err != nil {
		panic(err)
	}
	conn.ReadConfigFile("/etc/ceph/ceph.conf")
	err = conn.Connect()
	if err != nil {
		exitWithError("rados connection fail: %s", err)
	}
	return conn
}

func (s *CephImageSource) Matches(r *http.Request) bool {
	return r.Method == "GET" && r.URL.Query().Get("cns") != "" && r.URL.Query().Get("cid") != "" && s.Config.EnableCeph
}

func (s *CephImageSource) GetImage(req *http.Request) ([]byte, error) {
	co := parseObj(req)
	return s.fetchObject(co)
}

func parseObj(req *http.Request) CephObject {
	ns := req.URL.Query().Get("cns")
	id := req.URL.Query().Get("cid")
	return CephObject{ns, id}
}

func (s *CephImageSource) fetchObject(co CephObject) ([]byte, error) {
	ioctx, err := s.Connection.OpenIOContext(PoolName)
	defer ioctx.Destroy()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1048676)
	_, err = ioctx.GetXattr(co.NameSpace, co.OID, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func init() {
	RegisterSource(ImageSourceTypeCeph, NewCephImageSource)
}
