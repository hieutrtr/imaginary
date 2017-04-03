package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
)

var c *Ceph

func TestConnect(t *testing.T) {
	c = &Ceph{
		CephConfig: CephConfig{
			ConfigPath: "/etc/ceph/ceph.conf",
			Enable:     true,
			UseBlock:   false,
			BlockURL:   "",
		},
		CephObject: CephObject{
			Pool: "test",
		},
	}
	err := c.Connect()
	if err != nil {
		t.Fatal("Fail to connect Ceph")
	}
}

func TestOpenContext(t *testing.T) {
	err := c.OpenContext()
	if err != nil {
		t.Fatal("Fail to open Ceph conext: ", c.Pool)
	}
}

func TestBindRequest(t *testing.T) {
	vars := map[string]string{
		"service": "test",
		"oid":     "testobjstats",
		"attr":    "thumb",
	}
	c.BindObject(vars)
	if c.CephObject.Pool != vars["service"] || c.CephObject.OID != vars["oid"] || c.CephObject.Attr != vars["attr"] {
		t.Fatal("BindObject is fail")
	}
}

var imgTest = []byte("ImageOfTesting")

func TestSetAttr(t *testing.T) {
	if c.OnContext() {
		err := c.SetAttr(imgTest)
		if err != nil {
			t.Fatal("Set Attr to object is fail")
		}
	}
}

func TestGetAttr(t *testing.T) {
	if c.OnContext() {
		buf, err := c.GetAttr()
		if err != nil {
			t.Fatal("Get Attr from object is fail")
		}
		if bytes.Compare(buf, imgTest) != 0 {
			t.Fatal("Get Attr from object is different with original")
		}
	}
}

func TestGetStat(t *testing.T) {
	if c.OnContext() {
		stats, err := c.GetStat()
		if err != nil {
			t.Fatal("Getting Stats from ceph is fail on error", err.Error())
		}
		fmt.Println(stats.ModTime)
	}
}

func TestGetCachedAttr(t *testing.T) {
	url := "http://imaginaryct.com/7d6e3468f88cccbde9e94062650d632786dd54ea/ads/123/thumbnail?width=100"
	req := httptest.NewRequest("GET", url, nil)
	attr := getCacheAttr(req)
	if attr != "thumbnail_width=100" {
		t.Fatal("fail to get attr for caching", attr)
	}
}
