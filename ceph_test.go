package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c *Ceph

func TestConnect(t *testing.T) {
	c = &Ceph{
		CephConfig: CephConfig{
			ConfigPath: "/etc/ceph/ceph.conf",
			Enable:     true,
		},
	}
	err := c.Connect()
	if err != nil {
		t.Fatal("Fail to connect Ceph")
	}
}

func TestOpenContext(t *testing.T) {
	pool := "test"
	err := c.OpenContext(pool)
	if err != nil {
		t.Fatal("Fail to open Ceph conext: ", pool)
	}
}

var cephObj *CephObject

func TestBindObject(t *testing.T) {
	vars := map[string]string{
		"service": "test",
		"oid":     "testobjstats",
		"attr":    DATA,
	}
	assert := assert.New(t)
	cephObj = BindObject(vars)
	assert.Equal(cephObj.Pool, vars["service"], "service is fail bound")
	assert.Equal(cephObj.OID, vars["oid"], "service is fail oid")
	assert.Equal(cephObj.Attr, vars["attr"], "service is fail attr")
}

var imgTest = []byte("ImageOfTesting")

func TestSetAttr(t *testing.T) {
	pool := "test"
	if c.OnContext(pool) {
		err := c.SetAttr(cephObj, imgTest)
		if err != nil {
			t.Fatal("Set Attr to object is fail")
		}
	}
}

func TestGetAttr(t *testing.T) {
	pool := "test"
	if c.OnContext(pool) {
		buf, err := c.GetAttr(cephObj)
		if err != nil {
			t.Fatal("Get Attr from object is fail")
		}
		if bytes.Compare(buf, imgTest) != 0 {
			t.Fatal("Get Attr from object is different with original")
		}
	}
}

func TestGetStat(t *testing.T) {
	pool := "test"
	if c.OnContext(pool) {
		stats, err := c.GetStat(cephObj)
		if err != nil {
			t.Fatal("Getting Stats from ceph is fail on error", err.Error())
		}
		fmt.Println(stats.ModTime)
	}
}
