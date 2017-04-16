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

	vars := gorilla.Vars(req)
	if s.UseBlock {
		return ioutil.ReadFile(s.GetBlockPath(BindObject(vars)))
	}

	buf, err := s.GetAttr(BindObject(vars))
	// if s.Attr != DATA && buf == nil {
	// 	s.Attr = DATA
	// 	buf, err = s.fetchObject(req)
	// 	LoggerDebug.Println("Object need to be cached")
	// } else {
	// 	req.Header.Set("cached", s.Attr)
	req.Header.Set("cached", DATA)
	// }
	if err != nil {
		return nil, err
	}

	if stat, err := s.GetStat(BindObject(vars)); err == nil {
		req.Header.Set("Last-Modified", stat.ModTime.String())
	}
	return buf, err
}

func init() {
	RegisterSource(ImageSourceTypeCeph, NewCephImageSource)
}
