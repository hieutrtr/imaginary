package main

import (
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
			},
		},
	}
	if config.EnableCeph {
		MakeConnection(cis)
	}
	return cis
}

func (s *CephImageSource) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "GET" && vars["cpool"] != "" && vars["coid"] != ""
}

func (s *CephImageSource) GetImage(req *http.Request) ([]byte, error) {
	s.BindRequest(req)
	return s.fetchObject()
}

func (s *CephImageSource) fetchObject() ([]byte, error) {
	if !s.IsEnable() {
		return nil, NewError("ceph: service is not supported", Unsupported)
	}
	if s.Context == nil {
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
