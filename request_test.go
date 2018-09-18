package documentdb

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPartitionKey struct {
	Prop string `json:"prop"`
}

func (t *TestPartitionKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		NewProp string `json:"newProp"`
	}{NewProp: t.Prop})
}

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

func TestUpsertHeaders(t *testing.T) {
	r, _ := http.NewRequest("POST", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	_ = req.UpsertHeaders("YXJpZWwNCg==")

	assert := assert.New(t)
	assert.NotEqual(req.Header.Get(HEADER_AUTH), "")
	assert.NotEqual(req.Header.Get(HEADER_XDATE), "")
	assert.NotEqual(req.Header.Get(HEADER_VER), "")
	assert.Equal(req.Header.Get(HEADER_UPSERT), "true")
}

func TestPartitionKeyMarshalJSON(t *testing.T) {
	r, _ := http.NewRequest("GET", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	requestOptions := func(reqOpts *RequestOptions) {
		reqOpts.PartitionKey = &TestPartitionKey{"test"}
	}
	_ = req.RequestOptionsHeaders([]func(*RequestOptions){requestOptions})

	assert := assert.New(t)
	assert.Equal([]string{"{\"newProp\":\"test\"}"}, req.Header[HEADER_PARTITION_KEY])
}

func TestPartitionKeyAsInt(t *testing.T) {
	r, _ := http.NewRequest("GET", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	requestOptions := func(reqOpts *RequestOptions) {
		reqOpts.PartitionKey = 1
	}
	_ = req.RequestOptionsHeaders([]func(*RequestOptions){requestOptions})

	assert := assert.New(t)
	assert.Equal([]string{"[1]"}, req.Header[HEADER_PARTITION_KEY])
}

func TestPartitionKeyAsString(t *testing.T) {
	r, _ := http.NewRequest("GET", "link", &bytes.Buffer{})
	req := ResourceRequest("/dbs/b5NCAA==/", r)
	requestOptions := func(reqOpts *RequestOptions) {
		reqOpts.PartitionKey = "1"
	}
	_ = req.RequestOptionsHeaders([]func(*RequestOptions){requestOptions})

	assert := assert.New(t)
	assert.Equal([]string{"[\"1\"]"}, req.Header[HEADER_PARTITION_KEY])
}
