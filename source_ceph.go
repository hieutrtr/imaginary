package main

import (
	"io/ioutil"
	"net/http"

	gorilla "github.com/gorilla/mux"
)

// ImageSourceTypeCeph name of regiter source
const ImageSourceTypeCeph ImageSourceType = "ceph"

// CephImageSource source to register for Ceph
type CephImageSource struct {
	Ceph
}

// NewCephImageSource create new ceph image source
func NewCephImageSource(config *SourceConfig) ImageSource {
	cis := &CephImageSource{
		Ceph: Ceph{
			CephConfig: CephConfig{
				ConfigPath: config.CephConfig,
				Enable:     config.EnableCeph,
				UseBlock:   config.UseCephBlock,
				BlockURL:   config.CephBlockURL,
			},
		},
	}
	if config.EnableCeph && !config.UseCephBlock {
		err := MakeConnection(cis)
		if err != nil {
			exitWithError("Ceph connection was fail with config: %s", config.CephConfig)
		}
	}
	return cis
}

func (s *CephImageSource) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "GET" && vars["cpool"] != "" && vars["coid"] != ""
}

// GetImage from ceph
func (s *CephImageSource) GetImage(req *http.Request) ([]byte, error) {
	if !s.IsEnable() {
		return nil, NewError("ceph: service is not supported", Unsupported)
	}
	s.BindRequest(req)
	if s.UseBlock {
		return s.readFromBlock()
	}
	return s.fetchObject()
}

func (s *CephImageSource) readFromBlock() ([]byte, error) {
	return ioutil.ReadFile(s.GetBlockPath())
}

func (s *CephImageSource) fetchObject() ([]byte, error) {
	if !s.OnContext() {
		err := s.OpenContext()
		if err != nil {
			return nil, err
		}
	}
	return s.GetData()
}

func init() {
	RegisterSource(ImageSourceTypeCeph, NewCephImageSource)
}
