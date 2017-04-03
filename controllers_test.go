package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var conf = &SourceConfig{
	EnableCeph: true,
	CephConfig: "/etc/ceph/ceph.conf",
}

// MockCephImageSource create new ceph image source
func MockCephImageSource(config *SourceConfig) *CephImageSource {
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

func MockCephGetImage(req *http.Request) ([]byte, error) {
	s := MockCephImageSource(conf)
	s.Pool = "ads"
	s.OID = "1234"
	s.Attr = getCacheAttr(req)
	fmt.Println(s)
	img, err := s.fetchObject()
	if img != nil {
		req.Header.Set("cached", s.Attr)
	} else if s.Attr != DATA {
		s.Attr = DATA
		img, err = s.fetchObject()
	}
	fmt.Println(s)
	return img, err
}

func MockImgProcess(req *http.Request) string {
	if req.Header.Get("cached") != "" {
		return "ImgProcess Origin"
	}
	return "ImgProcess Normal"
}

func TestGetUncachedCephImageController(t *testing.T) {
	target := "http://10.60.6.12:8069/949161a8529db9f02ed7a1ac7772daa35e617d4f/ads/12345/thumbnail?width=100"
	req := httptest.NewRequest("GET", target, nil)

	img, err := MockCephGetImage(req)
	if err != nil {
		t.Fatal("Getting image from source fail with error", err)
	}

	fmt.Println(MockImgProcess(req))

	if imageRouting(req, img) != "cache" {
		t.Fatalf("Image is not cached with attribute <%s>\n", req.Header.Get("cached"))
	}
}

func TestUploadCephImageController(t *testing.T) {
	target := "http://10.60.6.12:8069/upload/ads/12345"
	img, _ := ioutil.ReadFile("./fixtures/imaginary.jpg")
	req := httptest.NewRequest("POST", target, nil)

	fmt.Println(MockImgProcess(req))

	if imageRouting(req, img) != "upload" {
		t.Fatal("Should not upload image")
	}
}
