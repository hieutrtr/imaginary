package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const formFieldName = "file"
const maxMemory int64 = 1024 * 1024 * 64

const ImageSourceTypeBody ImageSourceType = "payload"

type BodyImageSource struct {
	Config *SourceConfig
}

func NewBodyImageSource(config *SourceConfig) ImageSource {
	return &BodyImageSource{config}
}

func (s *BodyImageSource) Matches(r *http.Request) bool {
	return r.Method == "POST" || r.Method == "PUT"
}

func (s *BodyImageSource) GetImage(r *http.Request) ([]byte, error) {
	if isFormBody(r) {
		return readFormBody(r)
	}
	return readRawBody(r)
}

func (s *BodyImageSource) GetCache(r *http.Request) ([]byte, error) {
	return nil, nil
}

func (s *BodyImageSource) Delete(r *http.Request) error {
	return nil
}

func isFormBody(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/")
}

func readFormBody(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if len(buf) == 0 {
		err = ErrEmptyBody
	}

	return buf, err
}

func formField(r *http.Request) string {
	if field := r.URL.Query().Get("field"); field != "" {
		return field
	}
	return formFieldName
}

func readRawBody(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func init() {
	RegisterSource(ImageSourceTypeBody, NewBodyImageSource)
}
