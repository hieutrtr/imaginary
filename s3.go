package main

import (
	"bytes"
	"net/http"
	"os"
	"strconv"

	gorilla "github.com/gorilla/mux"
	s3 "github.com/minio/minio-go"
)

var (
	s3EndPoint        = os.Getenv("S3_ENDPOINT")
	s3AccessKeyID     = os.Getenv("S3_ACCESS_KEY_ID")
	s3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	s3UseSSL          = os.Getenv("S3_USE_SSL")
	s3ObjContentType  = "application/octet-stream"
)

// S3 : Struct of S3 handler
type S3 struct {
	Client  *s3.Client
	Service string
	OID     string
}

// S3Ready : check if missing required S3 env
func S3Ready() {
	if s3EndPoint == "" || s3AccessKeyID == "" || s3SecretAccessKey == "" || s3UseSSL == "" {
		exitWithError("need environment : S3_ENDPOINT, S3_ACCESS_KEY_ID, S3_SECRET_ACCESS_KEY, S3_USE_SSL")
	}
}

// Connect : connect s3 storage server
func (s *S3) Connect() error {
	var err error
	var useSSL bool
	if useSSL, err = strconv.ParseBool(s3UseSSL); err != nil {
		useSSL = false
	}
	s.Client, err = s3.New(s3EndPoint, s3AccessKeyID, s3SecretAccessKey, useSSL)
	return err
}

// BindRequest : Bind http request to get object
func (s *S3) BindRequest(req *http.Request) {
	vars := gorilla.Vars(req)
	s.Service = vars["service"]
	s.OID = vars["oid"]
}

// SetData push data to s3
func (s *S3) SetData(buf []byte) error {
	_, err := s.Client.PutObject(s.Service, s.OID, bytes.NewReader(buf), s3ObjContentType)
	return err
}

// GetData fetch object attribute DATA from ceph
func (s *S3) GetData() ([]byte, error) {
	obj, err := s.Client.GetObject(s.Service, s.OID)
	if err != nil {
		return nil, err
	}
	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, stat.Size)
	obj.Read(data)
	return data, nil
}
