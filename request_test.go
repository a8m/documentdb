package documentdb

import (
	"bytes"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestResourceRequest(t *testing.T) {
	assert := assert.New(t)
	req := ResourceRequest("/dbs/b5NCAA==/", &http.Request{})
	assert.Equal(req.rType, "dbs")
	assert.Equal(req.rId, "b5NCAA==")
}

func TestDefaultHeaders(t *testing.T) {
	r, _ := http.NewRequest("GET", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	_ = req.DefaultHeaders("YXJpZWwNCg==")

	assert := assert.New(t)
	assert.NotEqual(req.Header.Get(HEADER_AUTH), "")
	assert.NotEqual(req.Header.Get(HEADER_XDATE), "")
	assert.NotEqual(req.Header.Get(HEADER_VER), "")
}
