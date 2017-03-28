package main

import "testing"

func TestGetCacheAttr(t *testing.T) {
	url := "http://imaginary:8069/testpool/123/thumbnail"
	attr := getCacheAttr(url)
	if attr != "thumbnail" {
		t.Fatal("Wrong attribute is gotten in", url)
	}
	url = "http://imaginary:8069/testpool/123/wrongaction"
	attr = getCacheAttr(url)
	if attr != DATA {
		t.Fatal("Wrong attribute is gotten in", url)
	}
}
