package documentdb

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
//	"github.com/stretchr/testify/assert"
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

func (s *MockServer) AssertHeaders(headers ...string)

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
	s := ServerFactory(`{"_self": "Id"}`, 500)
	defer s.Close()
	client := &Client{}
	client.Url = s.URL
	client.Config = Config{"YXJpZWwNCg=="}
	db := &Database{}
	client.Read("/dbs/b5NCAA==/", db)
	defer s.Close()
	fmt.Println(db)
	fmt.Println(s.Header)
	fmt.Println(s.Body)
}
