package main

import (
	"fmt"
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
	if config.EnableCeph && !config.UseCephBlock {
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
		err := MakeConnection(cis)
		if err != nil {
			exitWithError("Ceph connection was fail %s with config: %s", fmt.Sprint(err), config.CephConfig)
		}
		return cis
	}
	return nil
}

func (s *CephImageSource) Matches(r *http.Request) bool {
	vars := gorilla.Vars(r)
	return r.Method == "GET" && vars["service"] != "" && vars["oid"] != ""
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
	cached := s.Attr
	buf, err := s.fetchObject()
	if s.Attr != DATA {
		if buf == nil {
			s.Attr = DATA
			buf, err = s.fetchObject()
			cached = ""
		}
	}

	if stat, err := s.GetStat(); err == nil {
		req.Header.Set("Last-Modified", stat.ModTime.String())
	}
	req.Header.Set("cached", cached)
	if cached == "" {
		LoggerDebug.Println("Object need to be cached")
	}
	return buf, err
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
	return s.GetAttr()
}

func init() {
	RegisterSource(ImageSourceTypeCeph, NewCephImageSource)
}
