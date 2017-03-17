package main

import (
	"fmt"
	"net/http"

	gorilla "github.com/gorilla/mux"
)

// S3ConnectionType name of regiter connection
const S3ConnectionType ConnectionType = "s3"

// S3Connection to register for S3
type S3Connection struct {
	S3
}

// NewS3Connection create new ceph connection
func NewS3Connection(config *ConnectionConfig) Connection {
	if config.EnableS3 {
		cc := &S3Connection{}
		err := MakeConnection(cc)
		if err != nil {
			exitWithError("S3 connection was fail %s", fmt.Sprint(err))
		}
		return cc
	}
	return nil
}

func (c *S3Connection) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "POST" && vars["service"] != "" && vars["oid"] != ""
}

// Execute purpose of openning connection
func (c *S3Connection) Execute(r *http.Request, buf []byte) error {
	c.BindRequest(r)
	return c.SetData(buf)
}

func init() {
	RegisterConnection(S3ConnectionType, NewS3Connection)
}
