package main

import (
	"bytes"
	"testing"
)

var sc *CephImageSource

func MockConnect() {
	sc = &CephImageSource{
		Ceph: Ceph{
			CephConfig: CephConfig{
				ConfigPath: "/etc/ceph/ceph.conf",
				Enable:     true,
				UseBlock:   false,
				BlockURL:   "",
			},
		},
	}
	sc.Connect()

	pool := "test"
	sc.OpenContext(pool)
}

var sCephObj *CephObject

func init() {
	MockConnect()
	vars := map[string]string{
		"service": "test",
		"oid":     "testobjstats",
		"attr":    DATA,
	}
	sCephObj = BindObject(vars)
}

var sImgTest = []byte("ImageOfTesting")

func TestSetImage(t *testing.T) {
	pool := "test"
	if sc.OnContext(pool) {
		err := sc.SetAttr(sCephObj, imgTest)
		if err != nil {
			t.Fatal("Set Attr to object is fail")
		}
	}
}

func TestGetImage(t *testing.T) {
	pool := "test"
	if sc.OnContext(pool) {
		buf, _, err := sc.IndependGetImage(sCephObj)
		if err != nil {
			t.Fatal("Get Attr from object is fail")
		}
		if bytes.Compare(buf, imgTest) != 0 {
			t.Fatal("Get Attr from object is different with original")
		}
	}
}