package documentdb

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

// I more interested in the request, instead of the response
type RequestRecorder struct {
	Header	http.Header
	Body	string
}

type MockServer struct {
	*httptest.Server
	RequestRecorder
}

func (s *MockServer) Record(r *http.Request) {
	s.Header = r.Header
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	s.Body = string(b)
}

func (s *MockServer) AssertHeaders(t *testing.T, headers ...string) {
	assert := assert.New(t)
	for _, k := range headers {
		assert.NotNil(s.Header[k])
	}
}

func ServerFactory(resp ...interface{}) *MockServer {
	s := &MockServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the last request
		s.Record(r)
		if v, ok := resp[0].(int); ok {
			err := fmt.Errorf("Status code %d", v)
			http.Error(w, err.Error(), v)
		} else {
			fmt.Fprintln(w, resp[0])
		}
		resp = resp[1:]
	}))
	return s
}

func TestRead(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &Client{Url:s.URL, Config:Config{"YXJpZWwNCg=="}}

	// First call
	var db Database
	err := client.Read("/dbs/b5NCAA==/", &db)
	s.AssertHeaders(t, HEADER_XDATE, HEADER_AUTH, HEADER_VER)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second Call
	// When StatusCode != StatusOK
	err = client.Read("/dbs/b5NCAA==/", &db)
}
