package documentdb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RequestRecorder struct {
	Header http.Header
	Body   string
}

type MockServer struct {
	*httptest.Server
	RequestRecorder
	Status interface{}
}

func (m *MockServer) SetStatus(status int) {
	m.Status = status
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
			err := fmt.Errorf(`{"code": "500", "message": "DocumentDB error"}`)
			http.Error(w, err.Error(), v)
		} else {
			if status, ok := s.Status.(int); ok {
				w.WriteHeader(status)
			}
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
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	var db Database
	_, err := client.Read("/dbs/b7NTAS==/", &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.Read("/dbs/b7NCAA==/colls/Ad352/", &db)
	assert.Equal(err.Error(), "500, DocumentDB error")
}

func TestReadWithUserAgent(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	testUserAgent := "test/user agent"
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}
	client.UserAgent = testUserAgent

	// First call
	var db Database
	_, err := client.Read("/dbs/b7NTAS==/", &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion, HeaderUserAgent)
	assert.Equal(s.Header.Get(HeaderUserAgent), testUserAgent)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")
}

func TestQuery(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, 500)
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	var db Database
	_, err := client.Query("dbs", &Query{Query: "SELECT * FROM ROOT r"}, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	s.AssertHeaders(t, HeaderContentLength, HeaderContentType, HeaderIsQuery)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.Read("/dbs/b7NCAA==/colls/Ad352/", &db)
	assert.Equal(err.Error(), "500, DocumentDB error")
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusCreated)
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	var db Database
	_, err := client.Create("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.Create("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, "9", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.Create("dbs", tDoc, &doc)
	assert.Equal(err.Error(), "500, DocumentDB error")
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`10`, 500)
	s.SetStatus(http.StatusNoContent)
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	_, err := client.Delete("/dbs/b7NTAS==/")
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Nil(err, "err should be nil")

	// Second Call, when StatusCode != StatusOK
	_, err = client.Delete("/dbs/b7NCAA==/colls/Ad352/")
	assert.Equal(err.Error(), "500, DocumentDB error")
}

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusOK)
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	var db Database
	_, err := client.Replace("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.Replace("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, "9", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.Replace("dbs", tDoc, &doc)
	assert.Equal(err.Error(), "500, DocumentDB error")
}

func TestExecute(t *testing.T) {
	assert := assert.New(t)
	s := ServerFactory(`{"_colls": "colls"}`, `{"id": "9"}`, 500)
	s.SetStatus(http.StatusOK)
	defer s.Close()
	client := &Client{Url: s.URL, Config: NewConfig(&Key{Key: "YXJpZWwNCg=="})}

	// First call
	var db Database
	_, err := client.Execute("dbs", `{"id": 3}`, &db)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(db.Colls, "colls", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Second call
	var doc, tDoc Document
	tDoc.Id = "9"
	_, err = client.Execute("dbs", tDoc, &doc)
	s.AssertHeaders(t, HeaderXDate, HeaderAuth, HeaderVersion)
	assert.Equal(doc.Id, "9", "Should fill the fields from response body")
	assert.Nil(err, "err should be nil")

	// Last Call, when StatusCode != StatusOK && StatusCreated
	_, err = client.Execute("dbs", tDoc, &doc)
	assert.Equal(err.Error(), "500, DocumentDB error")
}
