package documentdb

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type ClientStub struct {
	mock.Mock
}

func (c *ClientStub) Read(link string, ret interface{}) error {
	c.Called(link)
	return nil
}

func (c *ClientStub) Query(link, query string, ret interface{}) error {
	c.Called(link, query)
	return nil
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", Config{"config"})
	assert.IsType(client, &DocumentDB{}, "Should return DocumentDB object")
}

// TODO: Test failure
func TestReadDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadDatabase("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadCollection(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadCollection("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadDocument(t *testing.T) {
	type MyDocument struct {
		Document
		// Your external fields
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		IsAdmin bool   `json:"isAdmin,omitempty"`
	}
	var doc MyDocument
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link_doc").Return(nil)
	c.ReadDocument("self_link_doc", &doc)
	client.AssertCalled(t, "Read", "self_link_doc")
}

func TestReadStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadStoredProcedure("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadUserDefinedFunction("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "dbs").Return(nil)
	c.ReadDatabases()
	client.AssertCalled(t, "Read", "dbs")
}

func TestReadCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	dbLink := "dblink/"
	client.On("Read", dbLink + "colls/").Return(nil)
	c.ReadCollections(dbLink)
	client.AssertCalled(t, "Read", dbLink + "colls/")
}

func TestReadStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink + "sprocs/").Return(nil)
	c.ReadStoredProcedures(collLink)
	client.AssertCalled(t, "Read", collLink + "sprocs/")
}

func TestReadUserDefinedFunctions(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink + "udfs/").Return(nil)
	c.ReadUserDefinedFunctions(collLink)
	client.AssertCalled(t, "Read", collLink + "udfs/")
}

func TestReadDocuments(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink + "docs/").Return(nil)
	c.ReadDocuments(collLink, struct {}{})
	client.AssertCalled(t, "Read", collLink + "docs/")
}

func TestQueryDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "dbs", "SELECT * FROM ROOT r").Return(nil)
	c.QueryDatabases("SELECT * FROM ROOT r")
	client.AssertCalled(t, "Query", "dbs", "SELECT * FROM ROOT r")
}

func TestQueryCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "db_self_link/colls/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryCollections("db_self_link/", "SELECT * FROM ROOT r")
	client.AssertCalled(t, "Query", "db_self_link/colls/", "SELECT * FROM ROOT r")
}

func TestQueryStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "colls_self_link/sprocs/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryStoredProcedures("colls_self_link/", "SELECT * FROM ROOT r")
	client.AssertCalled(t, "Query", "colls_self_link/sprocs/", "SELECT * FROM ROOT r")
}
