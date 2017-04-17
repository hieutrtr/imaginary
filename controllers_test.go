package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var conf = &SourceConfig{
	EnableCeph: true,
	CephConfig: "/etc/ceph/ceph.conf",
}

var CtlImage = []byte("This is a little image")

func TestRoutingRequest(t *testing.T) {
	assert := assert.New(t)
	t.Run("Upload image from local to ceph", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload/test/abc", bytes.NewReader(CtlImage)) // Upload Request
		assert.Equal(UPLOAD, routingRequest(req), "Should be routed to uploaded request")
	})
	t.Run("Process origin image from ceph", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/abc", bytes.NewReader(CtlImage)) // Process Original Request
		assert.Equal(PROCESS, routingRequest(req), "Should be routed to process image request")
	})
	t.Run("Get cached image from ceph ", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test/abc/thumbnail?width=800", bytes.NewReader(CtlImage)) // Get Cache Request
		assert.Equal(CACHE, routingRequest(req), "Should be routed to get cache request")
	})
}
