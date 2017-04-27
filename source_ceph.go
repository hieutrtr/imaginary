package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	gorilla "github.com/gorilla/mux"
	"github.com/noahdesu/go-ceph/rados"
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
func (s *CephImageSource) IndependGetImage(co *CephObject) ([]byte, rados.ObjectStat, error) {

	var buf []byte
	var stat rados.ObjectStat
	var err error

	if !s.IsEnable() {
		return buf, stat, NewError("ceph: service is not supported", Unsupported)
	}

	if s.UseBlock {
		if buf, err = ioutil.ReadFile(s.GetBlockPath(co)); err != nil {
			return nil, stat, err
		}
		return buf, stat, nil
	}

	if buf, err = s.GetAttr(co); err != nil {
		return nil, stat, err
	}

	if stat, err = s.GetStat(co); err != nil {
		// req.Header.Set("Last-Modified", stat.ModTime.String())
		return nil, stat, err
	}

	return buf, stat, nil
}

// GetImage from ceph
func (s *CephImageSource) GetImage(req *http.Request) ([]byte, error) {
	if !s.IsEnable() {
		return nil, NewError("ceph: service is not supported", Unsupported)
	}
	fmt.Println("GetImage")

	vars := gorilla.Vars(req)
	if s.UseBlock {
		return ioutil.ReadFile(s.GetBlockPath(BindObject(vars)))
	}

	buf, err := s.GetAttr(BindObject(vars))
	if err != nil {
		if buf, err = ioutil.ReadFile("./fixtures/default_avatar.png"); err != nil {
			return nil, NewError(err.Error(), InternalError)
		}
	}

	if stat, err := s.GetStat(BindObject(vars)); err == nil {
		req.Header.Set("Last-Modified", stat.ModTime.String())
	}
	return buf, err
}

func init() {
	RegisterSource(ImageSourceTypeCeph, NewCephImageSource)
}
