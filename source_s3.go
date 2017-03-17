package main

import (
	"fmt"
	"net/http"

	gorilla "github.com/gorilla/mux"
)

// ImageSourceTypeS3 name of regiter source
const ImageSourceTypeS3 ImageSourceType = "s3"

// S3ImageSource source to register for S3
type S3ImageSource struct {
	S3
}

// NewS3ImageSource create new s3 image source
func NewS3ImageSource(config *SourceConfig) ImageSource {
	if config.EnableS3 {
		cis := &S3ImageSource{}
		err := MakeConnection(cis)
		if err != nil {
			exitWithError("S3 connection was fail %s with config: %s", fmt.Sprint(err))
		}
		return cis
	}
	return nil
}

func (s *S3ImageSource) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "GET" && vars["service"] != "" && vars["oid"] != ""
}

// GetImage from s3
func (s *S3ImageSource) GetImage(req *http.Request) ([]byte, error) {
	s.BindRequest(req)
	return s.GetData()
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSource)
}
